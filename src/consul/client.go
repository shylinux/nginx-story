package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

const (
	SERVICE  = "service"
	INSTANCE = "instance"
)

type client struct {
	ice.Hash
	short  string `data:"sess"`
	field  string `data:"time,sess,host,port"`
	create string `name:"create sess*=demo host port*=8500"`
	insert string `name:"insert name*=demo host=127.0.0.1 port*"`
	list   string `name:"list sess instance auto" help:"服务发现"`
}

func (s client) client(m *ice.Message, h string) *api.Client {
	return mdb.HashSelectTarget(m.Message, h, func(value ice.Maps) ice.Any {
		config := api.DefaultConfig()
		config.Address = fmt.Sprintf("%s:%s", value[tcp.HOST], value[tcp.PORT])
		if registry, err := api.NewClient(config); !m.Warn(err) {
			return registry
		}
		return nil
	}).(*api.Client)
}
func (s client) Inputs(m *ice.Message, arg ...string) {
	switch s.Hash.Inputs(m, arg...); arg[0] {
	case mdb.NAME:
		m.Push(arg[0], "mysql", "redis", "pulsar")
	case tcp.PORT:
		switch m.Option(mdb.NAME) {
		case "mysql":
			m.Push(arg[0], "10001", "3306")
		case "redis":
			m.Push(arg[0], "10002", "6379")
		case "pulsar":
			m.Push(arg[0], "10003", "6650")
		}
	}
}
func (s client) Insert(m *ice.Message, arg ...string) {
	m.Warn(s.client(m, m.Option(aaa.SESS)).Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name: m.Option(mdb.NAME), Address: m.Option(tcp.HOST), Port: kit.Int(m.Option(tcp.PORT)),
		ID: fmt.Sprintf("%s-%s-%s", m.Option(mdb.NAME), m.Option(tcp.HOST), m.Option(tcp.PORT)),
	}))
}
func (s client) Delete(m *ice.Message, arg ...string) {
	m.Warn(s.client(m, m.Option(aaa.SESS)).Agent().ServiceDeregister(m.Option(INSTANCE)))
}
func (s client) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 {
		s.Hash.List(m).Action(s.Hash.Create)
	} else if client := s.client(m, arg[0]); len(arg) == 1 {
		list, err := client.Agent().Services()
		if m.Warn(err) {
			return
		}
		for k, v := range list {
			m.Push(SERVICE, v.Service).Push(INSTANCE, k)
			m.Push(tcp.HOST, v.Address).Push(tcp.PORT, v.Port)
			m.Push("tags", kit.Join(v.Tags))
		}
		m.Sort(INSTANCE).PushAction(s.Delete).Action(s.Insert).StatusTimeCount()
	} else {
		service, _, _ := client.Agent().Service(arg[1], nil)
		m.Echo("%s", kit.Formats(service))
	}
}

func init() { ice.CodeCtxCmd(client{}) }
