package server

import (
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/misc/ssh"
	kit "shylinux.com/x/toolkits"
)

type cert struct {
	path string `data:"etc/conf/cert/"`
	list string `name:"list path auto upload" help:"证书"`
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
		}
	}
}
func (s cert) Show(m *ice.Message, arg ...string) {
	m.Cmdy(nfs.CAT, path.Join(m.Config(nfs.PATH), m.Option(nfs.PATH))).ProcessInner()
}

func init() { ice.CodeModCmd(cert{}) }
