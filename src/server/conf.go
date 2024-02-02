package server

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/nfs"
)

type conf struct{ ice.Lang }

func (s conf) Init(m *ice.Message, arg ...string) {
	s.Lang.Init(m, nfs.SCRIPT, m.Resource(""))
}
func init() { ice.Cmd("web.code.nginx.conf", conf{}) }
