package consul

import (
	"os"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	ice.Code
	ice.Hash
	darwin string `data:"https://releases.hashicorp.com/consul/1.15.2/consul_1.15.2_darwin_amd64.zip"`
	linux  string `data:"https://releases.hashicorp.com/consul/1.15.2/consul_1.15.2_linux_amd64.zip"`
	start  string `name:"start port*=8500"`
	list   string `name:"list path auto start install" help:"服务发现"`
}

func (s server) Start(m *ice.Message, arg ...string) {
	os.MkdirAll("usr/install/consul/bin", ice.MOD_DIR)
	os.Rename("usr/install/consul/consul", "usr/install/consul/bin/consul")
	s.Code.Daemon(m, "usr/install/consul/", "bin/consul", "agent", "-dev")
	kit.If(cli.IsSuccess(m.Message), func() { m.ProcessRefresh() })
}
func (s server) List(m *ice.Message, arg ...string) { s.Code.List(m, "consul", arg...) }

func init() { ice.CodeCtxCmd(server{}) }
