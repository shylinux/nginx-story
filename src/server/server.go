package server

import (
	"path"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	ice.Code

	source string `data:"http://mirrors.tencent.com/macports/distfiles/nginx/nginx-1.19.1.tar.gz"`
	start  string `name:"start port" help:"启动"`
	reload string `name:"reload" help:"重载"`
}

func (s server) Download(m *ice.Message, arg ...string) {
	s.Code.Download(m, m.Config(nfs.SOURCE), arg...)
}
func (s server) Build(m *ice.Message, arg ...string) {
	s.Code.Build(m, m.Config(nfs.SOURCE), "--with-http_ssl_module")
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
	s.Code.Start(m, m.Config(nfs.SOURCE), "sbin/nginx")
}
func (s server) Reload(m *ice.Message, arg ...string) {
	p := kit.Path(path.Join(m.Conf(cli.DAEMON, kit.Keym(nfs.PATH)), m.Option(tcp.PORT)))
	s.Code.System(m, p, "sbin/nginx", "-p", p, "-s", "reload")
}
func (s server) List(m *ice.Message, arg ...string) {
	if s.Code.List(m, m.Config(nfs.SOURCE), arg...); len(arg) == 0 || arg[0] == "" {
		m.Table(func(index int, value map[string]string, head []string) {
			m.PushAnchor(kit.Format("http://%s:%s", m.OptionUserWeb().Hostname(), value[tcp.PORT]))
		}).PushAction(s.Reload)
	}
}
func init() { ice.CodeModCmd(server{}) }
