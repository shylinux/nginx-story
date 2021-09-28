package server

import (
	"path"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/base/web"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	source string `data:"http://mirrors.tencent.com/macports/distfiles/nginx/nginx-1.19.1.tar.gz"`

	download string `name:"download" help:"下载"`
	build    string `name:"build" help:"构建"`
	start    string `name:"start" help:"启动"`
	reload   string `name:"start" help:"重载"`
	list     string `name:"list port path auto start build download" help:"服务器"`
}

func (s server) Download(m *ice.Message, arg ...string) {
	m.Cmdy(code.INSTALL, web.DOWNLOAD, m.Conf(tcp.SERVER, kit.META_SOURCE))
}
func (s server) Build(m *ice.Message, arg ...string) {
	m.Cmdy(code.INSTALL, cli.BUILD, m.Conf(tcp.SERVER, kit.META_SOURCE), "--with-http_ssl_module")
}
func (s server) Start(m *ice.Message, arg ...string) {
	m.Optionv(code.PREPARE, func(p string) []string {
		kit.Rewrite(path.Join(p, "conf/nginx.conf"), func(line string) string {
			if strings.HasPrefix(strings.TrimSpace(line), "listen") {
				return strings.ReplaceAll(line, kit.Split(line, "\t ", ";")[1], path.Base(p))
			}
			return line
		})
		return []string{"-p", kit.Path(p), "-g", "daemon off;"}
	})
	m.Cmdy(code.INSTALL, cli.START, m.Conf(tcp.SERVER, kit.META_SOURCE), "sbin/nginx")
}
func (s server) Reload(m *ice.Message, arg ...string) {
	p := m.Option(cli.CMD_DIR, kit.Path(path.Join(m.Conf(cli.DAEMON, kit.META_PATH), m.Option(tcp.PORT))))
	m.Cmdy(cli.SYSTEM, "sbin/nginx", "-p", p, "-s", "reload")
}
func (s server) List(m *ice.Message, arg ...string) {
	m.Cmdy(code.INSTALL, m.Conf(tcp.SERVER, kit.META_SOURCE), arg)
	if len(arg) == 0 || arg[0] == "" {
		m.Table(func(index int, value map[string]string, head []string) {
			u := m.OptionUserWeb()
			m.PushAnchor(kit.Format("http://%s:%s", u.Hostname(), value[tcp.PORT]))
			m.PushButton(cli.RELOAD)
		})
	}
}
func init() { ice.Cmd("web.code.nginx.server", server{}) }
