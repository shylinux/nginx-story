package server

import (
	"path"
	"path/filepath"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/lex"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/base/web"
	kit "shylinux.com/x/toolkits"
)

const (
	HTTP     = "http"
	CERT     = "cert"
	SERVER   = "server"
	LISTEN   = "listen"
	LOCATION = "location"
	UPSTREAM = "upstream"
	INCLUDE  = "include"

	SERVER_NAME      = "server_name"
	PROXY_PASS       = "proxy_pass"
	PROXY_SET_HEADER = "proxy_set_header"

	_CONF = ".conf"
)

type configs struct {
	path   string `data:"etc/conf/"`
	file   string `data:"nginx.conf"`
	create string `name:"create name* https=yes,no upstream* server*"`
	list   string `name:"list index path auto" help:"服务配置" icon:"nginx.png"`
}

func (s configs) Search(m *ice.Message, arg ...string) {
	if arg[0] == mdb.FOREACH {
		s.List(m.Spawn(ice.Maps{mdb.FIELDS: ""})).Table(func(value ice.Maps) {
			m.PushSearch(mdb.TYPE, web.LINK, mdb.NAME, value[mdb.NAME], mdb.TEXT, s.host(m, value[mdb.NAME], value[LISTEN]), value)
		})
	}
}
func (s configs) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case mdb.NAME:
		m.Cmd(nfs.DIR, path.Join(m.Config(nfs.PATH), SERVER), mdb.NAME, func(value ice.Maps) {
			m.Push(arg[0], strings.TrimSuffix(value[mdb.NAME], _CONF))
		})
	case CERT:
		m.Cmd(nfs.DIR, path.Join(m.Config(nfs.PATH), CERT), mdb.NAME, kit.Dict(nfs.DIR_REG, kit.ExtReg("key")), func(value ice.Maps) {
			m.Push(arg[0], strings.TrimSuffix(value[mdb.NAME], ".key"))
		})
	case UPSTREAM:
		m.Cmd(nfs.DIR, path.Join(m.Config(nfs.PATH), UPSTREAM), mdb.NAME, func(value ice.Maps) {
			m.Push(arg[0], strings.TrimSuffix(value[mdb.NAME], _CONF))
		})
	case SERVER:
		m.Cmdy(tcp.PORT, mdb.INPUTS, arg)
	}
}
func (s configs) Create(m *ice.Message, arg ...string) {
	m.Cmd(nfs.DEFS, path.Join(m.Config(nfs.PATH), SERVER, m.Option(mdb.NAME)+_CONF), nfs.Template(m.Message, kit.Select(SERVER+_CONF, "servers.conf", m.Option("https") == "yes")))
	m.Cmd(nfs.DEFS, path.Join(m.Config(nfs.PATH), LOCATION, m.Option(UPSTREAM)+_CONF), nfs.Template(m.Message, LOCATION+_CONF))
	m.Cmd(nfs.DEFS, path.Join(m.Config(nfs.PATH), UPSTREAM, m.Option(UPSTREAM)+_CONF), nfs.Template(m.Message, UPSTREAM+_CONF))
}
func (s configs) List(m *ice.Message, arg ...string) *ice.Message {
	stats := map[string]int{}
	conf, _ := s.parse(m, m.Config(nfs.PATH), m.Config(nfs.FILE), ice.Map{}, []string{}, stats)
	if len(arg) == 0 {
		list := map[string]bool{}
		m.Cmd(tcp.PORT, mdb.INPUTS, SERVER, func(value ice.Maps) { list[value[SERVER]] = true })
		kit.For(kit.Value(conf, kit.Keys(HTTP, SERVER)), func(index int, value ice.Map) {
			m.Push(mdb.INDEX, index).Push(mdb.NAME, kit.Format(value[SERVER_NAME])).Push(LISTEN, kit.Format(value[LISTEN]))
			p := kit.Format(kit.Value(value, kit.Keys(LOCATION, nfs.PS, PROXY_PASS)))
			server := kit.Format(kit.Value(conf, kit.Keys(HTTP, UPSTREAM, strings.TrimPrefix(p, "http://"), SERVER)))
			status := kit.Select("online", "offline", !list[server] && strings.HasPrefix(server, "127.0.0.1"))
			m.Push(PROXY_PASS, p).Push(SERVER, server).Push(mdb.STATUS, status)
			m.PushButton(s.Open)
			stats[status]++
		})
		m.Sort("status,name", "str_r", "str").Action(s.Create).StatusTime(stats)
		return m
	}
	server := kit.Value(conf, kit.Keys(HTTP, SERVER, arg[0])).(ice.Map)
	if p := s.host(m, kit.Format(server[SERVER_NAME]), kit.Format(server[LISTEN])); len(arg) == 1 {
		kit.For(kit.Value(server, LOCATION), func(path string, value ice.Any) {
			m.Push(nfs.PATH, path).Push(PROXY_PASS, kit.Value(value, PROXY_PASS))
		})
		m.StatusTimeCount(tcp.HOST, p)
	} else {
		m.EchoIFrame(p+arg[1]).StatusTimeCount(tcp.HOST, p+arg[1])
	}
	return m
}
func (s configs) Open(m *ice.Message, arg ...string) {
	m.ProcessOpen(s.host(m, m.Option(mdb.NAME), m.Option(LISTEN)))
}

func init() { ice.CodeModCmd(configs{}) }

func (s configs) host(m *ice.Message, host, port string) string {
	if strings.HasSuffix(port, " ssl") {
		if port == "443 ssl" {
			return kit.Format("https://%s", host)
		} else {
			return kit.Format("https://%s:%s", host, strings.TrimSuffix(port, " ssl"))
		}
	} else {
		if port == "80" {
			return kit.Format("http://%s", host)
		} else {
			return kit.Format("http://%s:%s", host, port)
		}
	}
}
func (s configs) parse(m *ice.Message, dir, file string, conf ice.Map, block []string, stats map[string]int) (ice.Map, []string) {
	m.Cmd(lex.SPLIT, dir+file, kit.Dict(lex.SPLIT_SPACE, "\t ;"), func(ls []string) {
		switch ls[0] {
		case INCLUDE:
			list, err := filepath.Glob(path.Join(dir, ls[1]))
			m.Warn(err)
			for _, file := range list {
				conf, block = s.parse(m, dir, strings.TrimPrefix(file, dir), conf, block, stats)
			}
		case HTTP, "events", "types":
			block = []string{ls[0]}
		case "}":
			block = kit.Slice(block, 0, -1)
		case PROXY_SET_HEADER:
			kit.Value(conf, kit.Keys(block, ls[0], ls[1]), strings.Join(ls[2:], lex.SP))
		case LOCATION:
			block = append(block, kit.Keys(LOCATION, ls[1]))
		case UPSTREAM:
			block = []string{HTTP, kit.Keys(UPSTREAM, ls[1])}
			stats[UPSTREAM]++
		case SERVER:
			if ls[1] == "{" {
				kit.Value(conf, kit.Keys(HTTP, SERVER, "-2"), kit.Dict())
				block = []string{HTTP, kit.Keys(SERVER, "-3")}
				stats[SERVER]++
			} else {
				kit.Value(conf, kit.Keys(block, SERVER), ls[1])
			}
			break

			m.Debug("what %v %v", ls, block)
			if len(block) == 1 || len(block) > 1 && strings.HasPrefix(block[1], SERVER+nfs.PT) {
				kit.Value(conf, kit.Keys(block, kit.Keys(SERVER, "-2")), kit.Dict())
				block = append(block, kit.Keys(SERVER, "-3"))
				stats[SERVER]++
				break
			}
			fallthrough
		default:
			m.Debug("what %v %v", ls, block)
			kit.Value(conf, kit.Keys(block, ls[0]), strings.Join(ls[1:], lex.SP))
		}
	})
	return conf, block
}
