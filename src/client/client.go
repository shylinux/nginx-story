package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/base/web"
	kit "shylinux.com/x/toolkits"
)

type client struct {
	ice.Hash
	short string `data:"sess"`
	field string `data:"time,sess,proto,host,port,path"`

	create string `name:"create sess=biz proto=http host=localhost port=10004 path=/" help:"创建"`
	list   string `name:"list sess@key auto create" help:"客户端"`
}

func (s client) List(m *ice.Message, arg ...string) {
	if s.Hash.List(m, arg...); len(arg) > 0 && arg[0] != "" {
		m.Cmdy(web.SPIDE_GET, kit.Format("%s://%s:%s/%s", m.Append(tcp.PROTO), m.Append(tcp.HOST), m.Append(tcp.PORT), m.Append(nfs.PATH)))
	}
}

func init() { ice.CodeModCmd(client{}) }
