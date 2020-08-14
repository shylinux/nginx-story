package server

import (
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/icebergs/base/gdb"
	"github.com/shylinux/icebergs/base/nfs"
	"github.com/shylinux/icebergs/base/tcp"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"

	"os"
	"path"
	"runtime"
	"strings"
)

const (
	NGINX  = "nginx"
	SERVER = "server"
)

var Index = &ice.Context{Name: NGINX, Help: "nginx",
	Configs: map[string]*ice.Config{
		SERVER: {Name: SERVER, Help: "服务器", Value: kit.Data(
			"windows", "https://nginx.org/download/nginx-1.8.1.zip",
			"darwin", "https://nginx.org/download/nginx-1.8.1.tar.gz",
			"linux", "https://nginx.org/download/nginx-1.8.1.tar.gz",
		)},
	},
	Commands: map[string]*ice.Command{
		ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},
		ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},

		SERVER: {Name: "server port=auto 查看:button=auto 返回:button 启动:button 编译:button 下载:button", Help: "服务器", Action: map[string]*ice.Action{
			"download": {Name: "download", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy("web.code.install", "download", m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)))
			}},
			"compile": {Name: "compile", Help: "编译", Hand: func(m *ice.Message, arg ...string) {
				name := path.Base(strings.TrimSuffix(strings.TrimSuffix(m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)), ".tar.gz"), "zip"))
				m.Option(cli.CMD_DIR, path.Join("usr/install", name))
				m.Cmdy(cli.SYSTEM, "./configure")
				m.Cmdy(cli.SYSTEM, "make")
			}},
			gdb.START: {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				if m.Option("port") == "" {
					m.Option("port", m.Cmdx(tcp.PORT, "get"))
				}
				p := kit.Format("var/daemon/%s", m.Option("port"))
				os.MkdirAll(path.Join(p, "logs"), ice.MOD_DIR)
				os.MkdirAll(path.Join(p, "bin"), ice.MOD_DIR)
				os.MkdirAll(p, ice.MOD_DIR)

				name := path.Base(strings.TrimSuffix(strings.TrimSuffix(m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)), ".tar.gz"), "zip"))
				m.Cmd(cli.SYSTEM, "cp", "-r", path.Join("usr/install", name, "conf"), p)
				m.Cmd(cli.SYSTEM, "cp", "-r", path.Join("usr/install", name, "html"), p)
				m.Cmd(cli.SYSTEM, "cp", "-r", path.Join("usr/install", name, "objs/nginx"), path.Join(p, "bin"))
				m.Cmd(cli.SYSTEM, "sed", "-i", kit.Format("s/80/%s/", m.Option("port")), path.Join(p, "conf/nginx.conf"))

				m.Option(cli.CMD_DIR, p)
				m.Cmdy(cli.DAEMON, "bin/nginx", "-p", kit.Path(p))
			}},
			gdb.STOP: {Name: "stop", Help: "停止", Hand: func(m *ice.Message, arg ...string) {
				p := kit.Format("var/daemon/%s", m.Option("port"))
				m.Option(cli.CMD_DIR, p)
				m.Cmdy(cli.SYSTEM, "bin/nginx", "-p", kit.Path(p), "-s", "stop")
			}},
			gdb.RELOAD: {Name: "reload", Help: "重载", Hand: func(m *ice.Message, arg ...string) {
				p := kit.Format("var/daemon/%s", m.Option("port"))
				m.Option(cli.CMD_DIR, p)
				m.Cmdy(cli.SYSTEM, "bin/nginx", "-p", kit.Path(p), "-s", "reload")
			}},
			gdb.RESTART: {Name: "restart", Help: "重启", Hand: func(m *ice.Message, arg ...string) {
				p := kit.Format("var/daemon/%s", m.Option("port"))
				m.Option(cli.CMD_DIR, p)
				m.Cmdy(cli.SYSTEM, "bin/nginx", "-p", kit.Path(p), "-s", "restart")
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if len(arg) == 0 {
				u := kit.ParseURL(m.Option(ice.MSG_USERWEB))
				m.Cmd(cli.DAEMON).Table(func(index int, value map[string]string, head []string) {
					if strings.HasPrefix(value["name"], "bin/nginx") {
						m.Push("time", value["time"])
						m.Push("port", path.Base(value["dir"]))
						m.Push("name", value["name"])
						m.Push("web", kit.Format("http://%s:%s", u.Hostname(), path.Base(value["dir"])))

					}
				})
				m.PushAction("重载", "重启", "停止")
				return
			}
			m.Cmdy(nfs.CAT, "var/daemon/"+arg[0]+"/conf/nginx.conf")
		}},
	},
}

func init() { code.Index.Register(Index, nil) }
