package studio

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/web"
)

type env struct {
	create string `name:"create origin* name* icons"`
	list   string `name:"list client.name auto"`
}

func (s env) Create(m *ice.Message, arg ...string) {
	m.Cmdy(web.SPIDE, m.ActionKey(), arg)
}
func (s env) Remove(m *ice.Message, arg ...string) {
	m.Cmdy(web.SPIDE, m.ActionKey(), arg)
}
func (s env) List(m *ice.Message, arg ...string) {
	if m.Cmdy(web.SPIDE, arg).PushAction(s.Remove); len(arg) == 0 {
		m.Action(s.Create)
	}
}

func init() { ice.CodeModCmd(env{}) }
