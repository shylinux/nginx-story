package server

import (
	"os"
	"path"
	"runtime"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/base/web"
	kit "shylinux.com/x/toolkits"
)

const (
	SBIN_NGINX = "./sbin/nginx"
)

type server struct {
	ice.Code
	source string `data:"http://mirrors.tencent.com/macports/distfiles/nginx/nginx-1.19.1.tar.gz"`
	action string `data:"test,error,reload,conf,make"`
	start  string `name:"start port*=10000" help:"启动"`
	test   string `name:"test path*=/" help:"测试"`
	error  string `name:"error" help:"日志"`
	reload string `name:"reload" help:"重载"`
	conf   string `name:"conf" help:"配置"`
	make   string `name:"make" help:"编译"`
	list   string `name:"list port path auto start build download" help:"服务器"`
}

func (s server) Inputs(m *ice.Message, arg ...string) {
	if arg[0] == nfs.PATH {
		s.System(m, path.Join(m.Option(nfs.DIR), "conf"), "grep", "-rh", "location")
		list := kit.Dict()
		for _, v := range strings.Split(m.Result(), ice.NL) {
			if strings.HasPrefix(strings.TrimSpace(v), "#") {
				continue
			}
			if strings.TrimSpace(v) == "" {
				continue
			}
			list[kit.Slice(kit.Split(v), -2)[0]] = struct{}{}
		}
		for _, k := range kit.SortedKey(list) {
			m.Push(arg[0], k)
		}
		return
	}
	s.Code.Inputs(m, arg...)
}
func (s server) Build(m *ice.Message, arg ...string) {
	args := []string{}
	kit.Fetch(m.Configv(source{}, nfs.MODULE), func(key string, value string) {
		args = append(args, kit.Format("--add-module=%s", kit.Path(value)))
	})
	if runtime.GOOS == cli.LINUX {
		s.Code.Build(m, "", "--with-http_ssl_module", "--with-http_v2_module", args)
	} else {
		s.Code.Build(m, "", "--with-http_v2_module", args)
	}
}
func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Start(m, "", SBIN_NGINX, func(p string) []string {
		os.MkdirAll(path.Join(p, "logs"), ice.MOD_DIR)
		nfs.Rewrite(m.Message, path.Join(p, "conf/nginx.conf"), func(line string) string {
			if strings.HasPrefix(strings.TrimSpace(line), "listen") {
				return strings.Replace(line, kit.Split(line, "\t ", ";")[1], path.Base(p), 1)
			}
			return line
		})
		return []string{"-p", kit.Path(p), "-g", "daemon off;"}
	})
}
func (s server) Stop(m *ice.Message, arg ...string) {
	s.Code.System(m, m.Option(nfs.DIR), SBIN_NGINX, "-p", nfs.PWD, "-s", "stop")
	m.Option(mdb.HASH, "")
	s.Code.Daemon(m.Spawn(), m.Option(nfs.DIR), cli.STOP)
	s.Code.ToastSuccess(m)
}
func (s server) Test(m *ice.Message, arg ...string) {
	m.EchoIFrame(kit.Format("http://%s:%s", web.UserWeb(m).Hostname(), m.Option(tcp.PORT)))
}
func (s server) Error(m *ice.Message, arg ...string) {
	m.Cmdy(nfs.CAT, path.Join(m.Option(cli.DIR), "logs/error.log"))
}
func (s server) Reload(m *ice.Message, arg ...string) {
	s.Code.System(m, m.Option(nfs.DIR), SBIN_NGINX, "-p", nfs.PWD, "-s", "reload")
}
func (s server) Conf(m *ice.Message, arg ...string) {
	s.Code.Field(m, ice.GetTypeKey(source{}), kit.Simple(m.Option(nfs.DIR)+ice.PS, "conf/nginx.conf", "43"), arg...)
}
func (s server) Make(m *ice.Message, arg ...string) {
	s.Code.ToastLong(m, "编译中...", m.Option(nfs.DIR))
	s.Stream(m, s.Path(m, ""), cli.MAKE, "-j8")
	s.Stream(m, s.Path(m, ""), cli.MAKE, "install")
	s.Stop(m)
	s.Code.ToastLong(m, "停止中...", m.Option(nfs.DIR))
	m.Sleep("3s")
	s.Code.ToastLong(m, "启动中...", m.Option(nfs.DIR))
	s.Start(m)
	m.Sleep("1s")
	s.Code.Toast(m, "启动成功", m.Option(nfs.DIR))
	m.ProcessRefresh()
}
func (s server) List(m *ice.Message, arg ...string) { s.Code.List(m, "", arg...) }

func init() { ice.CodeModCmd(server{}) }
