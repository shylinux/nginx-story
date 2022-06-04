package server

import (
	"path"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/nfs"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	ice.Code
	source string `data:"http://mirrors.tencent.com/macports/distfiles/nginx/nginx-1.19.1.tar.gz"`
	start  string `name:"start port" help:"启动"`
	reload string `name:"reload" help:"重载"`
	list   string `name:"list port path auto start build download" help:"代理"`
}

func (s server) Build(m *ice.Message, arg ...string) {
	s.Code.Build(m, "", "--with-http_ssl_module")
}
func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Start(m, "", "sbin/nginx", func(p string) []string {
		kit.Rewrite(path.Join(p, "conf/nginx.conf"), func(line string) string {
			if strings.HasPrefix(strings.TrimSpace(line), "listen") {
				return strings.ReplaceAll(line, kit.Split(line, "\t ", ";")[1], path.Base(p))
			}
			return line
		})
		return []string{"-p", "./", "-g", "daemon off;"}
	})
}
func (s server) Reload(m *ice.Message, arg ...string) {
	s.Code.System(m, m.Option(nfs.DIR), "sbin/nginx", "-p", "./", "-s", "reload")
}
func (s server) List(m *ice.Message, arg ...string) {
	if s.Code.List(m, "", arg...); len(arg) == 0 || arg[0] == "" {
		s.PushLink(m).PushAction(s.Reload)
	}
}
func init() { ice.CodeModCmd(server{}) }
