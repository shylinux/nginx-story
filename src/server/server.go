package server

import (
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"

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
		SERVER: {Name: "server port=auto path=auto auto 启动:button 构建:button 下载:button", Help: "服务器", Action: map[string]*ice.Action{
			"download": {Name: "download", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(code.INSTALL, "download", m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)))
			}},
			"build": {Name: "build", Help: "构建", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(code.INSTALL, "build", m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)))
			}},
			"start": {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				m.Optionv("prepare", func(p string) []string {
					kit.Rewrite(path.Join(p, "conf/nginx.conf"), func(line string) string {
						if strings.HasPrefix(strings.TrimSpace(line), "listen") {
							return strings.ReplaceAll(line, kit.Split(line, "\t ", ";")[1], path.Base(p))
						}
						return line
					})
					return []string{"-p", kit.Path(p)}
				})
				m.Cmdy(code.INSTALL, "start", m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)), "sbin/nginx")
			}},
			"reload": {Name: "reload", Help: "重载", Hand: func(m *ice.Message, arg ...string) {
				p := m.Option(cli.CMD_DIR, kit.Path(path.Join(m.Conf(cli.DAEMON, kit.META_PATH), m.Option(kit.SSH_PORT))))
				m.Cmdy(cli.SYSTEM, "sbin/nginx", "-p", p, "-s", "reload")
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Cmdy(code.INSTALL, path.Base(m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS))), arg)
			if len(arg) == 0 || arg[0] == "" {
				m.PushAction("重载")
			}
		}},
	},
}

func init() { code.Index.Register(Index, nil) }
