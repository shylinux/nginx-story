package server

import (
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/icebergs/base/gdb"
	"github.com/shylinux/icebergs/base/nfs"
	"github.com/shylinux/icebergs/base/tcp"
	"github.com/shylinux/icebergs/base/web"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"

	"net/http"
	"os"
	"path"
	"strings"
)

const (
	NGINX = "nginx"

	SERVER = "server"
	CLIENT = "client"
	BENCH  = "bench"
)

var Index = &ice.Context{Name: "nginx", Help: "nginx",
	Configs: map[string]*ice.Config{
		SERVER: {Name: SERVER, Help: "服务器", Value: kit.Data(
			"source", "http://nginx.org/download/nginx-1.8.1.tar.gz",
		)},
	},
	Commands: map[string]*ice.Command{
		ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},
		ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},

		SERVER: {Name: "server port 查看:button=auto 启动:button 编译:button 下载:button", Help: "服务器", Action: map[string]*ice.Action{
			"install": {Name: "install", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				// 下载
				source := m.Conf(SERVER, "meta.source")
				p := path.Join(m.Conf("web.code._install", "meta.path"), path.Base(source))
				if _, e := os.Stat(p); e != nil {
					msg := m.Cmd(web.SPIDE, "dev", web.CACHE, http.MethodGet, source)
					m.Cmd(web.CACHE, web.WATCH, msg.Append(web.DATA), p)
				}

				// 解压
				m.Option(cli.CMD_DIR, m.Conf("web.code._install", "meta.path"))
				m.Cmdy(cli.SYSTEM, "tar", "xvf", path.Base(source))
			}},
			"compile": {Name: "compile", Help: "编译", Hand: func(m *ice.Message, arg ...string) {
				// 编译
				source := m.Conf(SERVER, "meta.source")
				m.Option(cli.CMD_DIR, path.Join(m.Conf("web.code._install", "meta.path"), strings.TrimSuffix(path.Base(source), ".tar.gz")))
				m.Cmdy(cli.SYSTEM, "./configure")
				m.Cmdy(cli.SYSTEM, "make")

				// 链接
				m.Cmd(nfs.LINK, "bin/nginx", path.Join(m.Option(cli.CMD_DIR), "objs/nginx"))
			}},
			gdb.START: {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				if m.Option("port") == "" {
					m.Option("port", m.Cmdx(tcp.PORT, "get"))
				}

				p := kit.Format("var/daemon/%s", m.Option("port"))
				os.MkdirAll(path.Join(p, "logs"), ice.MOD_DIR)
				os.MkdirAll(p, ice.MOD_DIR)
				source := m.Conf(SERVER, "meta.source")
				m.Cmd(cli.SYSTEM, "cp", "-r", path.Join(m.Conf("web.code._install", "meta.path"), strings.TrimSuffix(path.Base(source), ".tar.gz"), "conf"), p)
				m.Cmd(cli.SYSTEM, "cp", "-r", path.Join(m.Conf("web.code._install", "meta.path"), strings.TrimSuffix(path.Base(source), ".tar.gz"), "html"), p)

				m.Cmd(cli.SYSTEM, "sed", "-i", kit.Format("s/80/%s/", m.Option("port")), path.Join(p, "conf/nginx.conf"))
				m.Cmdy(cli.DAEMON, "bin/nginx", "-p", kit.Path(p))
			}},
			gdb.STOP: {Name: "stop", Help: "停止", Hand: func(m *ice.Message, arg ...string) {
				p := kit.Format("var/daemon/%s", m.Option("port"))
				m.Cmdy(cli.SYSTEM, "bin/nginx", "-p", kit.Path(p), "-s", "stop")
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Split(m.Cmdx(cli.SYSTEM, "sh", "-c", "ps aux|grep nginx|grep -v grep"),
				"USER PID CPU MEM VSZ RSS TTY STAT START TIME COMMAND", " ", "\n")
			m.Table(func(index int, value map[string]string, head []string) {
				u := kit.ParseURL(m.Option(ice.MSG_USERWEB))
				ls := strings.Split(value["COMMAND"], " -p ")
				if len(ls) > 1 {
					if ls = kit.Split(ls[1], "/", "/"); len(ls) > 1 {
						m.Push("port", ls[len(ls)-1])
						m.Push("web", kit.Format("http://%s:%s", u.Hostname(), ls[len(ls)-1]))
						return
					}
				}
				m.Push("port", "80")
				m.Push("web", kit.Format("http://%s:%s", u.Hostname(), "80"))
			})
			m.Appendv(ice.MSG_APPEND, "USER", "PID", "STAT", "START", "port", "web", "COMMAND")
			m.PushAction("停止")
		}},
	},
}

func init() { code.Index.Register(Index, nil) }
