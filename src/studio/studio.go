package studio

import (
	"net/http"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/lex"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/web"
	kit "shylinux.com/x/toolkits"
)

const (
	ENV    = "env"
	URL    = "url"
	METHOD = "method"
	PARAMS = "params"
	HEADER = "header"
	COOKIE = "cookie"
)

type studio struct {
	ice.Hash
	field string `data:"time,hash,name,description,method,url,type,params,header,cookie"`
	list  string `name:"list hash env@key auto" help:"接口测试" icon:"studio.png"`
}

func (s studio) Inputs(m *ice.Message, arg ...string) {
	switch m.Option(ctx.ACTION) {
	case ENV:
		m.Cmdy(web.SPIDE).CutTo(web.CLIENT_NAME, arg[0])
	case PARAMS:
	case HEADER:
		switch arg[0] {
		case mdb.NAME:
			m.Cmdy(nfs.DIR, "", mdb.NAME, kit.Dict(nfs.DIR_ROOT, nfs.TemplatePath(m, HEADER)))
		case mdb.VALUE:
			m.Push(arg[0], strings.Split(m.Cmdx(nfs.CAT, m.Option(mdb.NAME), kit.Dict(nfs.DIR_ROOT, nfs.TemplatePath(m, HEADER))), lex.NL))
		}
	case COOKIE:
		switch arg[0] {
		case mdb.NAME:
			m.Push(arg[0], ice.MSG_SESSID)
		case mdb.VALUE:
			m.Push(arg[0], m.Option(ice.MSG_SESSID))
		}
	default:
		switch s.Hash.Inputs(m, arg...); arg[0] {
		case METHOD:
			m.Push(arg[0], http.MethodGet, http.MethodPut, http.MethodPost, http.MethodHead, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace)
		}
	}
}
func (s studio) Request(m *ice.Message, arg ...string) {
	args := []string{}
	kit.For(kit.UnMarshal(m.Option(PARAMS)), func(key string, value string) { args = append(args, key, value) })
	m.Options(web.SPIDE_HEADER, kit.UnMarshal(m.Option(HEADER)), web.SPIDE_COOKIE, kit.UnMarshal(m.Option(COOKIE)))
	m.Cmdy(web.SPIDE, m.OptionDefault(ENV, ice.DEV), web.SPIDE_RAW, m.Option(METHOD), m.Option(URL), m.Option(mdb.TYPE), args)
	m.Render(ice.RENDER_RAW)
}
func (s studio) Save(m *ice.Message, arg ...string) {
	s.Hash.Modify(m, m.OptionSimple(mdb.HASH, PARAMS, HEADER, COOKIE)...)
}
func (s studio) List(m *ice.Message, arg ...string) {
	s.Hash.List(m, kit.Slice(arg, 0, 1)...).Action(s.Create).PushAction(s.Remove).Display("")
}

func init() { ice.CodeModCmd(studio{}) }
