package server

import (
	"path"
	"path/filepath"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/lex"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

const (
	HTTP     = "http"
	SERVER   = "server"
	LISTEN   = "listen"
	LOCATION = "location"
	UPSTREAM = "upstream"
	INCLUDE  = "include"

	SERVER_NAME      = "server_name"
	PROXY_PASS       = "proxy_pass"
	PROXY_SET_HEADER = "proxy_set_header"
)

type configs struct {
	path string `data:"etc/conf/"`
	file string `data:"nginx.conf"`
	list string `name:"list index path auto" help:"服务配置"`
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
		case "events", HTTP, "types":
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
			if len(block) == 1 || len(block) > 1 && strings.HasPrefix(block[1], SERVER+nfs.PT) {
				kit.Value(conf, kit.Keys(block, kit.Keys(SERVER, "-2")), kit.Dict())
				block = append(block, kit.Keys(SERVER, "-3"))
				stats[SERVER]++
				break
			}
			fallthrough
		default:
			kit.Value(conf, kit.Keys(block, ls[0]), strings.Join(ls[1:], lex.SP))
		}
	})
	return conf, block
}
func (s configs) List(m *ice.Message, arg ...string) {
	stats := map[string]int{}
	conf, _ := s.parse(m, m.Config(nfs.PATH), m.Config(nfs.FILE), ice.Map{}, []string{}, stats)
	if len(arg) == 0 {
		kit.For(kit.Value(conf, kit.Keys(HTTP, SERVER)), func(index int, value ice.Map) {
			m.Push(mdb.INDEX, index).Push(mdb.NAME, kit.Format(value[SERVER_NAME]))
			m.Push(LISTEN, kit.Format(value[LISTEN]))
			p := kit.Format(kit.Value(value, kit.Keys(LOCATION, nfs.PS, PROXY_PASS)))
			m.Push(PROXY_PASS, p).Push(SERVER, kit.Format(kit.Value(conf, kit.Keys(HTTP, UPSTREAM, strings.TrimPrefix(p, "http://"), SERVER))))
		})
		m.Echo(kit.Formats(conf)).StatusTime(stats)
		ctx.DisplayStoryJSON(m)
		return
	}
	server := kit.Value(conf, kit.Keys(HTTP, SERVER, arg[0])).(ice.Map)
	p := kit.Format(server[LISTEN])
	if strings.HasSuffix(p, " ssl") {
		p = kit.Format("https://%s:%s", server[SERVER_NAME], strings.TrimSuffix(p, " ssl"))
	} else {
		p = kit.Format("http://%s:%s", server[SERVER_NAME], p)
	}
	if len(arg) == 1 {
		kit.For(kit.Value(server, LOCATION), func(path string, value ice.Any) {
			m.Push(nfs.PATH, path).Push(PROXY_PASS, kit.Value(value, PROXY_PASS))
		})
		m.StatusTimeCount(tcp.HOST, p)
	} else {
		m.EchoIFrame(p + arg[1])
		m.StatusTimeCount(tcp.HOST, p+arg[1])
	}
}
func init() { ice.CodeModCmd(configs{}) }
