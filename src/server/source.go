package server

import (
	"path"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type source struct {
	ice.Code
	make   string `name:"make" help:"编译"`
	ctags  string `name:"ctags" help:"索引"`
	module string `name:"module name=demo" help:"模块"`
	list   string `name:"list path=usr/install/nginx-1.19.1/src/ file=core/ngx_cycle.h@key line=39 auto" help:"源代码"`
}

func (s source) Make(m *ice.Message, arg ...string) {
	s.Stream(m, s.Path(m, m.Config(server{}, nfs.SOURCE)), cli.MAKE, "-j8")
	s.Stream(m, s.Path(m, m.Config(server{}, nfs.SOURCE)), cli.MAKE, "install")
}
func (s source) Ctags(m *ice.Message, arg ...string) {
	p := s.Path(m, m.Config(server{}, nfs.SOURCE), ice.SRC)
	s.System(m.Spawn(), "", kit.Simple("ctags", "-a", "-R", p, m.Configv(nfs.MODULE))...)
	s.System(m.Spawn(), p, "ctags", "-a", "-R")
	m.ProcessHold(ice.SUCCESS)
}
func (s source) Module(m *ice.Message, arg ...string) {
	p := m.Config(kit.Keys("module", m.Option(mdb.NAME)), kit.Format("src/%s/", m.Option(mdb.NAME)))
	s.Code.Module(m, p+kit.Format("ngx_http_%s_module.c", m.Option(mdb.NAME)), _module_template)
	s.Code.Module(m, p+"config", _module_config_template)
}
func (s source) Plugin(m *ice.Message, arg ...string) {
	if strings.HasPrefix(arg[1], "auto/") {
		m.Cmdy(mdb.PLUGIN, nfs.SH, arg[1])
		return
	}
	if arg[0] == "config" {
		m.Cmdy(mdb.PLUGIN, nfs.SH, arg[1])
		return
	}
	if arg[0] == "conf" {
		m.Echo(m.Config(mdb.PLUGIN))
		return
	}
	m.Cmdy(code.VIMER, mdb.PLUGIN, arg)
}
func (s source) List(m *ice.Message, arg ...string) {
	m.Cmdy(code.VIMER, arg).Option("modules", kit.Simple(s.Path(m, m.Config(server{}, nfs.SOURCE))+"/src/", kit.SortedValue(m.Configv(nfs.MODULE))))
	m.Option("plug", m.Config("show.plug"))
	m.Option("exts", m.Config("show.exts"))
	m.Option("tabs", m.Config("show.tabs"))
	if arg[0] != ctx.ACTION {
		m.Action(nfs.SAVE, s.Make, s.Ctags, s.Module)
	}
	if len(arg) > 1 && strings.HasPrefix(arg[1], "auto/") {
		m.Cmdy(nfs.CAT, path.Join(arg[0], arg[1]))
	}
	if len(arg) > 1 && arg[1] == "config" {
		m.Cmdy(nfs.CAT, path.Join(arg[0], arg[1]))
	}
}

func init() {
	ice.CodeModCmd(source{}, mdb.PLUGIN, kit.Dict(
		code.PREFIX, kit.Dict(
			"#", code.COMMENT,
		),
		code.KEYWORD, kit.Dict(
			"events", code.KEYWORD,
			"http", code.KEYWORD,
			"upstream", code.KEYWORD,
			"server", code.KEYWORD,
			"location", code.KEYWORD,

			"worker_processes", code.FUNCTION,
			"worker_connections", code.FUNCTION,
			"include", code.FUNCTION,
			"default_type", code.FUNCTION,
			"listen", code.FUNCTION,
			"server_name", code.FUNCTION,
			"error_page", code.FUNCTION,
			"root", code.FUNCTION,
			"index", code.FUNCTION,
		),
	))
}

const _module_config_template = `
ngx_addon_name=ngx_http_{{.Option "name"}}_module
NGX_ADDON_SRCS="$NGX_ADDON_SRCS $ngx_addon_dir/${ngx_addon_name}.c"
HTTP_MODULES="$HTTP_MODULES ${ngx_addon_name}"
`
const _module_template = `
#include <ngx_config.h>
#include <ngx_core.h>
#include <ngx_http.h>

ngx_module_t ngx_http_{{.Option "name"}}_module;

typedef struct {
	ngx_str_t echo;
} ngx_http_{{.Option "name"}}_loc_conf_t;

ngx_int_t
ngx_http_{{.Option "name"}}_handler(ngx_http_request_t *r) {
	ngx_int_t rc = ngx_http_discard_request_body(r);
	if (rc != NGX_OK) {
		return rc;
	}

	ngx_http_{{.Option "name"}}_loc_conf_t *dlcf = ngx_http_get_module_loc_conf(r, ngx_http_{{.Option "name"}}_module);
	ngx_str_t echo = dlcf->echo;

	r->headers_out.status = NGX_HTTP_OK;
	r->headers_out.content_length_n = echo.len;
	rc = ngx_http_send_header(r);
	if (rc != NGX_OK) {
		return rc;
	}

	ngx_buf_t *buf = ngx_create_temp_buf(r->pool, echo.len);
	ngx_memcpy(buf->pos, echo.data, echo.len);
	buf->last = buf->pos+echo.len;
	buf->last_buf = 1;
	ngx_log_error(NGX_LOG_ERR, r->connection->log, 0, "what %d", buf->pos);
	ngx_log_error(NGX_LOG_ERR, r->connection->log, 0, "what %d", buf->last);

	ngx_chain_t out = {buf, NULL};
	return ngx_http_output_filter(r, &out);
}

static char*
ngx_http_{{.Option "name"}}(ngx_conf_t *cf, ngx_command_t *cmd, void *conf) {
	ngx_http_core_loc_conf_t *clcf = ngx_http_conf_get_module_loc_conf(cf, ngx_http_core_module);
	clcf->handler = ngx_http_{{.Option "name"}}_handler;

	ngx_str_t *value = cf->args->elts;
	ngx_http_{{.Option "name"}}_loc_conf_t *dlcf = ngx_http_conf_get_module_loc_conf(cf, ngx_http_{{.Option "name"}}_module);
	dlcf->echo = value[1];
	return NGX_CONF_OK;
}
void *
ngx_http_{{.Option "name"}}_create_loc_conf(ngx_conf_t *cf) {
	return ngx_palloc(cf->pool, sizeof(ngx_http_{{.Option "name"}}_loc_conf_t));
}

static ngx_command_t ngx_http_{{.Option "name"}}_commands[] = {
	{
		ngx_string("{{.Option "name"}}"),
		NGX_HTTP_MAIN_CONF|NGX_HTTP_SRV_CONF|NGX_HTTP_LOC_CONF|NGX_CONF_1MORE,
		ngx_http_{{.Option "name"}},
		NGX_HTTP_LOC_CONF_OFFSET,
		0, NULL,
	},
	ngx_null_command
};

static ngx_http_module_t ngx_http_{{.Option "name"}}_module_ctx = {
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	&ngx_http_{{.Option "name"}}_create_loc_conf,
	NULL
};

ngx_module_t ngx_http_{{.Option "name"}}_module = {
	NGX_MODULE_V1,
	&ngx_http_{{.Option "name"}}_module_ctx,
	ngx_http_{{.Option "name"}}_commands,
	NGX_HTTP_MODULE,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	NGX_MODULE_V1_PADDING
};
`
