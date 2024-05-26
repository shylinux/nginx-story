package studio

import (
	"net/http"
	"path"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/lex"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/web"
	"shylinux.com/x/icebergs/base/web/html"
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
	ice.Hash
	tools  string `data:"env"`
	export string `data:"true"`
	field  string `data:"time,hash,name,method,url,params,header,cookie,auth,config"`
	create string `name:"create url* method* name*"`
	list   string `name:"studio env@key list" icon:"studio.png"`
}

func (s studio) Init(m *ice.Message, arg ...string) {
	web.AddPortalProduct(m.Message, "API Studio", `
一款网页版的接口测试工作台，用来进行接口测试。
`, 10.0)
}
func (s studio) Inputs(m *ice.Message, arg ...string) {
	switch m.Option(ctx.ACTION) {
	case ENV:
		m.Cmdy(web.SPIDE).CutTo(web.CLIENT_NAME, arg[0])
	case CONFIG:
		if arg[0] == mdb.NAME {
			defer m.Push(arg[0], html.PROFILE, html.DISPLAY, ctx.INDEX, ctx.ARGS)
		} else if arg[0] == mdb.VALUE {
			switch m.Option(mdb.NAME) {
			case html.PROFILE, html.DISPLAY:
				m.Cmd(nfs.DIR, "", mdb.NAME, kit.Dict(nfs.DIR_ROOT, nfs.TemplatePath(m.Message, CONFIG, m.Option(mdb.NAME)))).Table(func(value ice.Maps) { m.Push(arg[0], value[mdb.NAME]) })
				m.Push(arg[0], kit.ExtChange(path.Base(m.Option(URL)), kit.Select(nfs.HTML, nfs.JS, m.Option(mdb.NAME) == html.DISPLAY)))
			case ctx.INDEX, ctx.ARGS:
				s.Hash.Inputs(m, m.Option(mdb.NAME))
			}
		}
		fallthrough
	case PARAMS, HEADER, COOKIE, AUTH:
		switch arg[0] {
		case mdb.NAME:
			m.Cmdy(nfs.DIR, "", mdb.NAME, kit.Dict(nfs.DIR_ROOT, nfs.TemplatePath(m.Message, m.Option(ctx.ACTION))))
		case mdb.VALUE:
			m.Push(arg[0], kit.Filters(strings.Split(m.Cmdx(nfs.CAT, m.Option(mdb.NAME), kit.Dict(nfs.DIR_ROOT, nfs.TemplatePath(m.Message, m.Option(ctx.ACTION)))), lex.NL), ""))
		}
	default:
		switch s.Hash.Inputs(m, arg...); arg[0] {
		case METHOD:
			m.Push(arg[0], http.MethodGet, http.MethodPut, http.MethodPost, http.MethodHead, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace)
		}
	}
}
func (s studio) List(m *ice.Message, arg ...string) {
	s.Hash.List(m).PushAction(s.Remove).Action(s.Create).Display("").DisplayCSS("")
}
func (s studio) Request(m *ice.Message, arg ...string) {
	args, header := []string{}, kit.UnMarshal(m.Option(HEADER))
	if m.Option(METHOD) != http.MethodGet {
		if strings.HasPrefix(kit.Format(kit.Value(header, html.ContentType)), html.ApplicationJSON) {
			args = append(args, web.SPIDE_JSON)
		} else {
			args = append(args, web.SPIDE_FORM)
		}
	}
	kit.For(kit.UnMarshal(m.Option(PARAMS)), func(k, v string) { args = append(args, k, v) })
	kit.For(kit.UnMarshal(m.Option(AUTH)), func(k, v string) { kit.Value(header, html.Authorization, k+lex.SP+v) })
	m.Options(web.SPIDE_HEADER, header, web.SPIDE_COOKIE, kit.UnMarshal(m.Option(COOKIE)))
	m.Cmdy(web.SPIDE, m.OptionDefault(ENV, ice.OPS), web.SPIDE_DETAIL, m.Option(METHOD), m.Option(URL), args).Render(ice.RENDER_RAW)
}

func init() { ice.CodeModCmd(studio{}) }
