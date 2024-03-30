package server

import (
	"os"
	"path"
	"runtime"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

const (
	SBIN_NGINX      = "sbin/nginx"
	CONF_NGINX_CONF = "conf/nginx.conf"
	LOGS_ERROR_LOG  = "logs/error.log"
)

type server struct {
	ice.Code
	source string `data:"http://mirrors.tencent.com/macports/distfiles/nginx/nginx-1.19.1.tar.gz"`
	action string `data:"reload,stop,conf,test,error"`
	start  string `name:"start port*=10000" help:"启动"`
	reload string `name:"reload" help:"重载" icon:"bi bi-bootstrap-reboot"`
	conf   string `name:"conf" help:"配置"`
	test   string `name:"test path*=/" help:"测试" icon:"bi bi-clipboard-check"`
	error  string `name:"error" help:"日志" icon:"bi bi-calendar-week"`
	list   string `name:"list port path auto start build download" help:"服务器"`
}

func (s server) Init(m *ice.Message, arg ...string) {
	code.PackageCreate(m.Message, nfs.SOURCE, "nginx", "", "", s.Link(m))
}
func (s server) Build(m *ice.Message, arg ...string) {
	args := []string{}
	if runtime.GOOS == cli.LINUX {
		s.Code.Build(m, "", "--with-http_v2_module", "--with-http_ssl_module", "--with-http_auth_request_module", args)
	} else {
		s.Code.Build(m, "", "--with-http_v2_module", "--without-http_rewrite_module", args)
	}
	m.Cmdy(nfs.DIR, path.Join(s.Path(m, ""), "_install/sbin/nginx"))
}
func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Start(m, "", SBIN_NGINX, func(p string) []string {
		os.MkdirAll(path.Join(p, "logs"), ice.MOD_DIR)
		nfs.Rewrite(m.Message, path.Join(p, CONF_NGINX_CONF), func(line string) string {
			if strings.HasPrefix(strings.TrimSpace(line), LISTEN) {
				return strings.Replace(line, kit.Split(line, "\t ", ";")[1], path.Base(p), 1)
			}
			return line
		})
		return []string{"-p", kit.Path(p), "-g", "daemon off;"}
	})
}
func (s server) Reload(m *ice.Message, arg ...string) { s.cmds(m, arg...) }
func (s server) Stop(m *ice.Message, arg ...string)   { s.cmds(m, arg...); s.Code.Stop(m, arg...) }
func (s server) Conf(m *ice.Message, arg ...string) {
	m.Cmdy(nfs.CAT, path.Join(m.Option(nfs.DIR), CONF_NGINX_CONF)).ProcessInner()
}
func (s server) Test(m *ice.Message, arg ...string) {
	m.EchoIFrame(kit.Format("http://%s:%s", m.UserWeb().Hostname(), m.Option(tcp.PORT))).ProcessInner()
}
func (s server) Error(m *ice.Message, arg ...string) {
	m.Cmdy(nfs.CAT, path.Join(m.Option(cli.DIR), LOGS_ERROR_LOG))
}
func (s server) List(m *ice.Message, arg ...string) { s.Code.List(m, "", arg...) }

func init() { ice.CodeModCmd(server{}) }

func (s server) cmds(m *ice.Message, arg ...string) {
	defer m.ToastProcess()()
	p := m.OptionDefault(nfs.DIR, path.Join(ice.USR_LOCAL_DAEMON, m.Option(tcp.PORT)))
	s.Code.System(m, p, SBIN_NGINX, "-p", kit.Path(p), "-s", m.ActionKey())
}
