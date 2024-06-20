package server

import (
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/web"
	"shylinux.com/x/icebergs/misc/ssh"
	kit "shylinux.com/x/toolkits"
)

type cert struct {
	path    string `data:"etc/conf/cert/"`
	list    string `name:"list path auto"`
	deliver string `name:"deliver server*"`
}

func (s cert) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case web.SERVER:
		m.Cmd(web.SPACE, ice.OPS, web.SPACE).Table(func(value ice.Maps) {
			if value[mdb.TYPE] == web.SERVER {
				m.Push(arg[0], value[mdb.NAME])
			}
		})
	}
}
func (s cert) Deliver(m *ice.Message, arg ...string) {
	defer m.ToastProcess()()
	to := kit.Keys(ice.OPS, m.Option(web.SERVER), m.Option(ice.MSG_USERPOD))
	m.Cmd(nfs.DIR, "etc/conf/cert/").Table(func(value ice.Maps) {
		text := m.Cmdx(nfs.CAT, value[nfs.PATH])
		m.Cmd(web.SPACE, to, nfs.SAVE, value[nfs.PATH], text+"\n")
	})
}
func (s cert) Upload(m *ice.Message, arg ...string) {
	p := m.UploadSave(m.Config(nfs.PATH))
	m.Option(nfs.FILE, path.Base(kit.TrimExt(p, ssh.KEY, ssh.PEM)))
	m.Cmd(nfs.SAVE, path.Join(path.Dir(p), "cert.conf"), m.Template("cert.conf"))
}
func (s cert) Trash(m *ice.Message, arg ...string) {
	m.Trash(path.Join(m.Config(nfs.PATH), m.Option(nfs.PATH)))
}
func (s cert) List(m *ice.Message, arg ...string) {
	if m.Options(nfs.DIR_ROOT, m.Config(nfs.PATH)).Cmdy(nfs.CAT, arg); len(arg) == 0 {
		pem, key := false, false
		m.Table(func(value ice.Maps) {
			kit.If(kit.Ext(value[nfs.PATH]) == ssh.PEM, func() { pem = true })
			kit.If(kit.Ext(value[nfs.PATH]) == ssh.KEY, func() { key = true })
		})
		if !pem {
			m.EchoInfoButton("please upload cert pem", s.Upload)
		} else if !key {
			m.EchoInfoButton("please upload cert key", s.Upload)
		} else {
			m.Action(s.Upload, s.Deliver)
		}
	}
}
func (s cert) Show(m *ice.Message, arg ...string) {
	m.ProcessFloat(nfs.CAT, path.Join(m.Config(nfs.PATH), m.Option(nfs.PATH)), arg...)
}

func init() { ice.CodeModCmd(cert{}) }
