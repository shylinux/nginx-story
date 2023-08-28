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
	CONFIG = "config"
	AUTH   = "auth"
)

type studio struct {
	ice.Code
	ice.Hash
	field  string `data:"time,hash,name,description,method,url,type,params,header,cookie,auth,config"`
	create string `name:"create name* description method* url*"`
	list   string `name:"list env@key list" help:"接口测试" icon:"studio.png"`
}

func (s studio) Init(m *ice.Message, arg ...string) { s.Hash.Import(m) }
func (s studio) Exit(m *ice.Message, arg ...string) { s.Hash.Export(m) }
func (s studio) Inputs(m *ice.Message, arg ...string) {
	switch m.Option(ctx.ACTION) {
	case ENV:
		m.Cmdy(web.SPIDE).CutTo(web.CLIENT_NAME, arg[0])
	case CONFIG:
		if arg[0] == mdb.VALUE && m.Option(mdb.NAME) == "display" {
			m.Cmd(nfs.DIR, "", mdb.NAME, kit.Dict(nfs.DIR_ROOT, nfs.TemplatePath(m, CONFIG, m.Option(mdb.NAME)))).Table(func(value ice.Maps) {
				m.Push(arg[0], kit.MergeURL("/require/"+nfs.TemplatePath(m, CONFIG, m.Option(mdb.NAME), value[mdb.NAME]), ice.POD, m.Option(ice.MSG_USERPOD)))
			})
			break
		}
		fallthrough
	case PARAMS, HEADER, COOKIE, "auth":
		switch arg[0] {
		case mdb.NAME:
			m.Cmdy(nfs.DIR, "", mdb.NAME, kit.Dict(nfs.DIR_ROOT, nfs.TemplatePath(m, m.Option(ctx.ACTION))))
		case mdb.VALUE:
			m.Push(arg[0], strings.Split(m.Cmdx(nfs.CAT, m.Option(mdb.NAME), kit.Dict(nfs.DIR_ROOT, nfs.TemplatePath(m, m.Option(ctx.ACTION)))), lex.NL))
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
	header := kit.UnMarshal(m.Option(HEADER))
	kit.For(kit.UnMarshal(m.Option(AUTH)), func(key, value string) { kit.Value(header, web.Authorization, key+lex.SP+value) })
	m.Options(web.SPIDE_HEADER, header, web.SPIDE_COOKIE, kit.UnMarshal(m.Option(COOKIE)))
	m.Cmdy(web.SPIDE, m.OptionDefault(ENV, ice.DEV), web.SPIDE_RAW, m.Option(METHOD), m.Option(URL), m.Option(mdb.TYPE), args)
	m.Render(ice.RENDER_RAW)
}
func (s studio) Save(m *ice.Message, arg ...string) {
	s.Hash.Modify(m, m.OptionSimple(mdb.HASH, PARAMS, HEADER, COOKIE, AUTH, CONFIG)...)
}
func (s studio) List(m *ice.Message, arg ...string) {
	s.Hash.List(m).Action(s.Create).PushAction(s.Remove).Display("")
	m.StatusTimeCount(ENV, kit.Select(ice.DEV, arg, 0))
}

func init() { ice.CodeModCmd(studio{}) }
