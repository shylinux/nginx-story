package server

import (
	"path"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	ice.Code

	source string `data:"http://mirrors.tencent.com/macports/distfiles/nginx/nginx-1.19.1.tar.gz"`
	reload string `name:"start" help:"重载"`
}

func (s server) Download(m *ice.Message, arg ...string) {
	s.Code.Download(m, m.Config(cli.SOURCE), arg...)
}
func (s server) Build(m *ice.Message, arg ...string) {
	s.Code.Build(m, m.Config(cli.SOURCE), "--with-http_ssl_module")
}
func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Prepare(m, func(p string) []string {
		kit.Rewrite(path.Join(p, "conf/nginx.conf"), func(line string) string {
			if strings.HasPrefix(strings.TrimSpace(line), "listen") {
				return strings.ReplaceAll(line, kit.Split(line, "\t ", ";")[1], path.Base(p))
			}
			return line
		})
		return []string{"-p", kit.Path(p), "-g", "daemon off;"}
	})
	s.Code.Start(m, m.Config(cli.SOURCE), "sbin/nginx")
}
func (s server) Reload(m *ice.Message, arg ...string) {
	p := kit.Path(path.Join(m.Conf(cli.DAEMON, kit.META_PATH), m.Option(tcp.PORT)))
	s.Code.System(m, p, "sbin/nginx", "-p", p, "-s", "reload")
}
func (s server) List(m *ice.Message, arg ...string) {
	s.Code.List(m, m.Config(cli.SOURCE), arg...)
	if len(arg) == 0 || arg[0] == "" {
		u := m.OptionUserWeb()
		m.Table(func(index int, value map[string]string, head []string) {
			m.PushAnchor(kit.Format("http://%s:%s", u.Hostname(), value[tcp.PORT]))
		})
		m.PushAction(s.Reload)
	}
}
func init() { ice.Cmd("web.code.nginx.server", server{}) }
