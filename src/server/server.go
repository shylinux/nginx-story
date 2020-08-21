package server

import (
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/icebergs/base/gdb"
	"github.com/shylinux/icebergs/base/mdb"
	"github.com/shylinux/icebergs/base/tcp"
	"github.com/shylinux/icebergs/base/web"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"

	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
)

func _nginx_port(m *ice.Message, p string) {
	f, e := os.Open(path.Join(p, "conf/nginx.conf"))
	m.Assert(e)
	defer f.Close()

	b, e := ioutil.ReadAll(f)
	m.Assert(e)
	bio := bufio.NewScanner(bytes.NewBuffer(b))

	o, _, e := kit.Create(path.Join(p, "conf/nginx.conf"))
	m.Assert(e)
	defer o.Close()
	for bio.Scan() {
		if strings.HasPrefix(strings.TrimSpace(bio.Text()), "listen") {
			o.WriteString(kit.Format("        listen        %s;", path.Base(p)))
			o.WriteString("\n")
			continue
		}
		o.WriteString(bio.Text())
		o.WriteString("\n")
	}
}

const (
	NGINX  = "nginx"
	SERVER = "server"
)

var Index = &ice.Context{Name: NGINX, Help: "nginx",
	Configs: map[string]*ice.Config{
		NGINX: {Name: NGINX, Help: "服务器", Value: kit.Data(
			"windows", "https://nginx.org/download/nginx-1.8.1.zip",
			"darwin", "https://nginx.org/download/nginx-1.8.1.tar.gz",
			"linux", "https://nginx.org/download/nginx-1.8.1.tar.gz",
		)},
	},
	Commands: map[string]*ice.Command{
		ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},
		ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},

		SERVER: {Name: "server port=auto auto 启动:button 编译:button 下载:button", Help: "服务器", Action: map[string]*ice.Action{
			"download": {Name: "download", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(code.INSTALL, "download", m.Conf(NGINX, kit.Keys(kit.MDB_META, runtime.GOOS)))
			}},
			"compile": {Name: "compile", Help: "编译", Hand: func(m *ice.Message, arg ...string) {
				name := path.Base(strings.TrimSuffix(strings.TrimSuffix(m.Conf(NGINX, kit.Keys(kit.MDB_META, runtime.GOOS)), ".tar.gz"), "zip"))
				m.Option(cli.CMD_DIR, path.Join(m.Conf(code.INSTALL, kit.META_PATH), name))
				m.Cmdy(cli.SYSTEM, "./configure", "--prefix=./install")
				m.Cmdy(cli.SYSTEM, "make", "-j8")
				m.Cmdy(cli.SYSTEM, "make", "install")
			}},
			"start": {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				// 分配
				port, p := "", ""
				for {
					port = m.Cmdx(tcp.PORT, "select", port)
					p = path.Join(m.Conf(cli.DAEMON, kit.META_PATH), port)
					if _, e := os.Stat(p); e != nil && os.IsNotExist(e) {
						break
					}
					port = kit.Format(kit.Int(port) + 1)
				}

				// 复制
				name := path.Base(strings.TrimSuffix(strings.TrimSuffix(m.Conf(NGINX, kit.Keys(kit.MDB_META, runtime.GOOS)), ".tar.gz"), "zip"))
				from := kit.Path(path.Join(m.Conf(code.INSTALL, kit.META_PATH), name, "install"))
				m.Cmdy(cli.SYSTEM, "cp", "-r", from, p)
				_nginx_port(m, p)

				// 启动
				m.Option(cli.CMD_DIR, p)
				m.Cmdy(cli.DAEMON, "sbin/nginx", "-p", kit.Path(p))
			}},
			gdb.STOP: {Name: "stop", Help: "停止", Hand: func(m *ice.Message, arg ...string) {
				p := path.Join(m.Conf(cli.DAEMON, kit.META_PATH), m.Option(kit.MDB_PORT))
				m.Option(cli.CMD_DIR, p)
				m.Cmdy(cli.SYSTEM, "sbin/nginx", "-p", kit.Path(p), "-s", "stop")
			}},
			gdb.RELOAD: {Name: "reload", Help: "重载", Hand: func(m *ice.Message, arg ...string) {
				p := path.Join(m.Conf(cli.DAEMON, kit.META_PATH), m.Option(kit.MDB_PORT))
				m.Option(cli.CMD_DIR, p)
				m.Cmdy(cli.SYSTEM, "sbin/nginx", "-p", kit.Path(p), "-s", "reload")
			}},
			gdb.RESTART: {Name: "restart", Help: "重启", Hand: func(m *ice.Message, arg ...string) {
				p := path.Join(m.Conf(cli.DAEMON, kit.META_PATH), m.Option(kit.MDB_PORT))
				m.Option(cli.CMD_DIR, p)
				m.Cmdy(cli.SYSTEM, "sbin/nginx", "-p", kit.Path(p), "-s", "restart")
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if len(arg) > 0 && arg[0] != "" {
				m.Cmdy(mdb.RENDER, web.RENDER.Frame, "http://shylinux.com:"+m.Option("port"))
				return
			}

			u := kit.ParseURL(m.Option(ice.MSG_USERWEB))
			m.Cmd(cli.DAEMON).Table(func(index int, value map[string]string, head []string) {
				if strings.HasPrefix(value[kit.MDB_NAME], "sbin/nginx") {
					m.Push(kit.MDB_TIME, value[kit.MDB_TIME])
					m.Push(kit.MDB_PORT, path.Base(value[kit.MDB_DIR]))
					m.Push(kit.MDB_NAME, value[kit.MDB_NAME])
					m.Push(kit.MDB_LINK, m.Cmdx(mdb.RENDER, web.RENDER.A,
						kit.Format("http://%s:%s", u.Hostname(), path.Base(value[kit.MDB_DIR]))))
				}
			})
			m.PushAction("重载", "重启", "停止")
			m.Sort("time", "time_r")
		}},
	},
}

func init() { code.Index.Register(Index, nil) }
