package server

import (
	"path"
	"runtime"
	"strings"

	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/icebergs/base/gdb"
	"github.com/shylinux/icebergs/base/tcp"
	"github.com/shylinux/icebergs/base/web"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"
)

const (
	SERVER = "server"
)
const NGINX = "nginx"

var Index = &ice.Context{Name: NGINX, Help: "nginx",
	Configs: map[string]*ice.Config{
		SERVER: {Name: SERVER, Help: "服务器", Value: kit.Data(
			cli.WINDOWS, "https://nginx.org/download/nginx-1.8.1.zip",
			cli.DARWIN, "https://nginx.org/download/nginx-1.8.1.tar.gz",
			cli.LINUX, "https://nginx.org/download/nginx-1.8.1.tar.gz",
		)},
	},
	Commands: map[string]*ice.Command{
		SERVER: {Name: "server port path auto start build download", Help: "服务器", Action: map[string]*ice.Action{
			web.DOWNLOAD: {Name: "download", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(code.INSTALL, web.DOWNLOAD, m.Conf(SERVER, kit.Keym(runtime.GOOS)))
			}},
			gdb.BUILD: {Name: "build", Help: "构建", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(code.INSTALL, gdb.BUILD, m.Conf(SERVER, kit.Keym(runtime.GOOS)), "--with-http_ssl_module")
			}},
			gdb.START: {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				m.Optionv(code.PREPARE, func(p string) []string {
					kit.Rewrite(path.Join(p, "conf/nginx.conf"), func(line string) string {
						if strings.HasPrefix(strings.TrimSpace(line), "listen") {
							return strings.ReplaceAll(line, kit.Split(line, "\t ", ";")[1], path.Base(p))
						}
						return line
					})
					return []string{"-p", kit.Path(p), "-g", "daemon off;"}
				})
				m.Cmdy(code.INSTALL, gdb.START, m.Conf(SERVER, kit.Keym(runtime.GOOS)), "sbin/nginx")
			}},
			gdb.RELOAD: {Name: "reload", Help: "重载", Hand: func(m *ice.Message, arg ...string) {
				p := m.Option(cli.CMD_DIR, kit.Path(path.Join(m.Conf(cli.DAEMON, kit.META_PATH), m.Option(tcp.PORT))))
				m.Cmdy(cli.SYSTEM, "sbin/nginx", "-p", p, "-s", "reload")
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Cmdy(code.INSTALL, m.Conf(SERVER, kit.Keym(runtime.GOOS)), arg)
			if len(arg) == 0 || arg[0] == "" {
				m.Table(func(index int, value map[string]string, head []string) {
					u := kit.ParseURL(m.Option(ice.MSG_USERWEB))
					m.PushAnchor(kit.Format("http://%s:%s", u.Hostname(), value[tcp.PORT]))
					m.PushButton(gdb.RELOAD)
				})
			}
		}},
	},
}

func init() { code.Index.Register(Index, nil) }
