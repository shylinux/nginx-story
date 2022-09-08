package server

import (
	"path"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

const (
	SBIN_NGINX = "./sbin/nginx"
)

type server struct {
	ice.Code
	source string `data:"http://mirrors.tencent.com/macports/distfiles/nginx/nginx-1.19.1.tar.gz"`
	action string `data:"reload"`

	module string `name:"module name" help:"模块"module `
	reload string `name:"reload" help:"重载"`
	ctags  string `name:"ctags" help:"索引"`
	list   string `name:"list port path auto start build ctags module download" help:"代理"`
}

func (s server) Init(m *ice.Message, arg ...string) {
	m.Config("module", kit.List())
}

func (s server) Ctags(m *ice.Message, arg ...string) {
	m.Cmd(cli.SYSTEM, "ctags", "-a", "-R", path.Join(m.Cmdx(code.INSTALL, nfs.PATH, m.Config(nfs.SOURCE)), ice.SRC))
	m.Cmdy(nfs.DIR, "tags")
}
func (s server) Module(m *ice.Message, arg ...string) {
	m.Config("module.-2", kit.Format("src/%s/", m.Option("name")))
	m.Cmd(nfs.DEFS, kit.Format("src/%s/config", m.Option("name")), kit.Renders(`
ngx_addon_name=ngx_http_{{.Option "name"}}_module
NGX_ADDON_SRCS="$NGX_ADDON_SRCS $ngx_addon_dir/${ngx_addon_name}.c"
HTTP_MODULES="$HTTP_MODULES $ngx_addon_name"
`, m))
	m.Cmd(nfs.DEFS, kit.Format("src/%s/ngx_http_%s_module.c", m.Option("name"), m.Option("name")), kit.Renders(`
#include <ngx_config.h>
#include <ngx_core.h>
#include <ngx_http.h>

ngx_int_t
ngx_http_{{.Option "name"}}_handler(ngx_http_request_t *r) {
	return NGX_DONE;
}

static char*
ngx_http_{{.Option "name"}}(ngx_conf_t *cf, ngx_command_t *cmd, void *conf) {
	ngx_http_core_loc_conf_t *clcf = ngx_http_conf_get_module_loc_conf(cf, ngx_http_core_module);
	clcf->handler = ngx_http_{{.Option "name"}}_handler;
	return NGX_CONF_OK;
}

static ngx_command_t ngx_http_{{.Option "name"}}_commands[] = {
	{
		ngx_string("{{.Option "name"}}"),
		NGX_HTTP_MAIN_CONF|NGX_HTTP_SRV_CONF|NGX_HTTP_LOC_CONF|NGX_CONF_NOARGS,
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
	NULL,
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
`, m))
}
func (s server) Build(m *ice.Message, arg ...string) {
	args := []string{}
	kit.Fetch(m.Configv("module"), func(index int, value string) {
		args = append(args, kit.Format("--add-module=%s", kit.Path(value)))
	})
	s.Code.Build(m, "", "--with-http_ssl_module", args)
}
func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Start(m, "", SBIN_NGINX, func(p string) []string {
		kit.Rewrite(path.Join(p, "conf/nginx.conf"), func(line string) string {
			if strings.HasPrefix(strings.TrimSpace(line), "listen") {
				return strings.ReplaceAll(line, kit.Split(line, "\t ", ";")[1], path.Base(p))
			}
			return line
		})
		return []string{"-p", kit.Path(p), "-g", "daemon off;"}
	})
}
func (s server) Stop(m *ice.Message, arg ...string) {
	s.Code.System(m, m.Option(nfs.DIR), SBIN_NGINX, "-p", nfs.PWD, "-s", "stop")
}
func (s server) Reload(m *ice.Message, arg ...string) {
	s.Code.System(m, m.Option(nfs.DIR), SBIN_NGINX, "-p", nfs.PWD, "-s", "reload")
}
func (s server) List(m *ice.Message, arg ...string) {
	s.Code.List(m, "", arg...)
	m.Tables(func(value ice.Maps) {
		switch value[mdb.STATUS] {
		case "start":
			m.PushButton(s.Reload, s.Stop)
		default:
			m.PushButton(s.Start, s.Trash)
		}
	})
}
func init() { ice.CodeModCmd(server{}) }
