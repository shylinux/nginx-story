(function() { const REQUEST = "request", RESPONSE = "response"
const METHOD = "method", PARAMS = "params", STATUS = "status", HEADER = "header", COOKIE = "cookie"
const PLUGIN_STORY_EDITOR = "/plugin/story/editor.js"
const PLUGIN_STORY_MONACO = "/plugin/story/monaco.js"
const PLUGIN_STORY_JSON = "/plugin/story/json.js"
Volcanos(chat.ONIMPORT, {
	_init: function(can, msg, cb) { can.onmotion.clear(can), can.ui = can.onappend.layout(can), cb && cb(msg)
		can.db._list = {}, can.db._hash = can.misc.SearchHash(can), can.onimport._project(can, msg)
	},
	_project: function(can, msg) {
		var _select; msg.Table(function(value) {
			if (value.method == web.DELETE) {
				value.nick = can.page.Format(html.SPAN, "DEL", [METHOD, "DEL"])+lex.SP+can.page.Format(html.SPAN, value.name, mdb.NAME)
			} else {
				value.nick = can.page.Format(html.SPAN, value.method, [METHOD, value.method])+lex.SP+can.page.Format(html.SPAN, value.name, mdb.NAME)
			}
			var _target = can.onimport.item(can, value, function(event) { if (value._tabs) { return value._tabs.click() }
				value._tabs = can.onimport.tabs(can, [value], function(event) { can.db.current = value, can.misc.SearchHash(can, value.hash)
					if (can.onmotion.cache(can, function(save, load) { save({msg: can.db.msg, action: can.ui.action, request: can.ui.request, response: can.ui.response, plugin: can.ui.content._plugin})
						return load(value.hash, function(bak) { can.db.msg = bak.msg, can.ui.action = bak.action, can.ui.request = bak.request, can.ui.response = bak.response, can.ui.content._plugin = bak.plugin })
					}, can.ui.content)) { _target.click() } else {
						can.onimport._content(can, value), can.ui.toggle = can.onappend.toggle(can, can.ui.content)
						can.onimport._profile(can, value, can.page.Append(can, can.ui.profile.parentNode, [html.PROFILE])._target)
						can.onimport._display(can, value, can.page.Append(can, can.ui.display.parentNode, [html.DISPLAY])._target)
						can.onmotion.hidden(can, value._profile), can.onmotion.hidden(can, value._display)
					}
					can.ui.profile = value._profile, can.ui.display = value._display
					can.onmotion.select(can, can.ui.profile.parentNode, html.DIV_PROFILE, can.ui.profile)
					can.onmotion.select(can, can.ui.display.parentNode, html.DIV_DISPLAY, can.ui.display)
					can.onimport.layout(can)
				}, function(event) {
					delete(value._tabs)
					delete(can.ui.content._cache[value.hash])
					can.page.Remove(can, value._profile), can.page.Remove(can, value._display)
					can.onmotion.delay(can, function() { can.onexport.tabs(can) })
				}), can.onexport.tabs(can)
			}, null, can.ui.project); can.db._list[value.hash] = _target, (!_select || value.hash == can.db._hash[0]) && (_select = _target)
		})
		can.core.Next(can.misc.sessionStorage(can, [can.ConfIndex(), html.TABS]), function(hash, next) {
			var tabs = can.db._list[hash]; tabs && tabs.click(), can.onmotion.delay(can, next, 50)
		}, function() { _select && _select.click() })
	},
	_content: function(can, value) { can.db.msg = can.request()
		can.ui.action = can.page.Append(can, can.ui.content, [
			{view: html.ACTION, _init: function(target) { can.onappend._action(can, [
				{type: html.SELECT, name: METHOD, value: value.method, values: [web.GET, web.PUT, web.POST, web.DELETE]},
				{type: html.TEXT, name: web.URL, value: value.url, action: "key"},
				{type: html.BUTTON, value: REQUEST, style: html.NOTICE},
			], target) }},
		])._target
		can.ui.request = can.onimport._part(can, REQUEST, can.core.List([PARAMS, HEADER, COOKIE, aaa.AUTH, ctx.CONFIG], function(key) {
			return {name: key, show: function(event, target) { can.onimport._table(can, target, value, key) }}
		}), can.ui.content, [{type: html.BUTTON, value: nfs.SAVE}])
		can.ui.response = can.onimport._part(can, RESPONSE, [{name: mdb.DATA, show: function(event, target) {
			can.onimport._plugin(can, target)
		}}].concat(can.core.List([STATUS, HEADER, COOKIE], function(key) {
			return {name: key, show: function(event, target) { can.onimport._response(can, target, can.db.msg, key) }}
		})), can.ui.content, [{type: html.BUTTON, value: REQUEST, style: html.NOTICE}])
	},
	_profile: function(can, value, target) { value._profile = target
		return
		can.onappend.plugin(can, {
			index: web.CODE_VIMER, args: can.onexport._args(can, html.PROFILE), display: PLUGIN_STORY_EDITOR, style: html.OUTPUT,
			height: can.ConfHeight(), width: (can.ConfWidth()-can.ui.project.offsetWidth)/2,
		}, function(sub) { target._plugin = sub }, target)
	},
	_display: function(can, value, target) { value._display = target
		return
		can.onappend.plugin(can, {
			index: web.CODE_VIMER, args: can.onexport._args(can, html.DISPLAY), display: PLUGIN_STORY_MONACO, style: html.OUTPUT,
		}, function(sub) { target._plugin = sub }, target)
	},
	_part: function(can, name, list, target, action) {
		var ui = can.page.Append(can, target, [
			{view: name, list: [{view: html.TITLE, list: [{text: can.base.capital(name)}], _init: function(target) {
				can.onappend._action(can, action, target, null, true)
			}}, {view: html.ACTION, list: can.core.List(list, function(item) {
				return {view: [[html.ITEM, item.name], "", item.name], onclick: function(event) {
					can.onmotion.select(can, ui.action, "", event.currentTarget)
					if (can.onmotion.cache(can, function() { return item.name}, ui.content)) { return }
					item.show(event, ui.content)
				}, _init: function(target) { target._item = item }}
			})}, {view: html.CONTENT},
			]},
		]); can.page.SelectOne(can, ui.action, "", function(target) { target.click() }); return ui
	},
	_table: function(can, target, value, _key, keys) { var msg = can.request(); keys = keys || [mdb.NAME, mdb.VALUE]
		can.core.Item(can.base.Obj(value[_key]), function(key, value) { msg.Push(mdb.NAME, key), msg.Push(mdb.VALUE, value) })
		can.core.List(keys, function(key) { msg.Push(key, "") })
		function add(val, key) { return {type: html.TD, _init: function(target) {
			can.onappend._action(can, [{type: html.TEXT, name: key, value: val, _init: function(target) {
				can.onappend.figure(can, {name: key, run: function(event, cmds, cb) { var msg = can.request(event, {action: _key}, value)
					can.page.Select(can, target.parentNode.parentNode.parentNode, html.INPUT, function(target, index) { msg.Option(keys[index], target.value) })
					can.run(event, [ctx.ACTION, mdb.INPUTS, key], cb)
				}, _enter: function() {
					can.page.Append(can, table, [{type: html.TR, list: can.core.List(keys, function(key) { return add("", key) }) }])
				}}, target, function(sub, value) {})
			}}], target)
		}, list: []} } var table = can.onappend.table(can, msg, add, target); return table
	},
	_response: function(can, target, msg, type) { var _msg = can.request()
		msg.Table(function(value) { if (value.type == type) { _msg.Push(mdb.NAME, value.name), _msg.Push(mdb.VALUE, value.value) } })
		can.onappend.style(can, type, can.onappend.table(can, _msg, null, target))
	},
	_plugin: function(can, target) {
		can.db.msg.Table(function(value) { if (value.name == http.ContentType) {
			var config = can.onexport.config({}, can)
			if (config.display) {
				if (can.base.beginWith(config.profile, nfs.PS, web.HTTP)) {

				} else if (can.base.beginWith(config.profile, nfs.SRC, nfs.USR)) {

				} else {
					config.display = can.misc.Template(can, "/config/display/", config.display)
				}
				config.display = can.misc.Resource(can, config.display)
				config.display = can.base.MergeURL(config.display, "_vv", new Date().getTime())
				can.onappend.plugin(can, {
					index: config.index, args: config.args, msg: can.db.msg, display: config.display, style: config.style,
					height: can.ConfHeight()/2-1,
				}, function(sub) { can.ui.content._plugin = sub
					can.onmotion.hidden(can, sub._legend), can.onmotion.hidden(can, sub._option)
				}, target)
				return
			}
			switch (can.core.Split(value.value, ";")[0]) {
				case mime.ApplicationJSON: can.onappend._output(can, can.db.msg, PLUGIN_STORY_JSON, null, target, null, false); break
				default: can.onappend.table(can, can.db.msg, null, target), can.onappend.board(can, can.db.msg, target)
			}
		} }) 
	},
	layout: function(can) { can.page.style(can, can.ui.display, html.HEIGHT, can.ConfHeight()/2)
		can.page.style(can, can.ui.profile, html.HEIGHT, can.ConfHeight()-can.ui.display.offsetHeight, html.WIDTH, (can.ConfWidth()-can.ui.project.offsetWidth)/2)
		can.ui.layout(can.ConfHeight(), can.ConfWidth()), can.ui.toggle && can.ui.toggle.layout()
		can.ui.profile._plugin && can.ui.profile._plugin.onimport.size(can.ui.profile._plugin, can.ConfHeight()-can.ui.display.offsetHeight, (can.ConfWidth()-can.ui.project.offsetWidth)/2-1)
		can.ui.display._plugin && can.ui.display._plugin.onimport.size(can.ui.display._plugin, can.ConfHeight()/2-1, can.ConfWidth()-can.ui.project.offsetWidth)
		can.ui.content._plugin && can.ui.content._plugin.onimport.size(can.ui.content._plugin, can.ConfHeight()/2-1, can.ConfWidth()-can.ui.project.offsetWidth-can.ui.profile.offsetWidth)
	},
}, [""])
Volcanos(chat.ONACTION, {
	save: function(event, can, button) { can.runAction(can.onexport.request(event, can), button, [], function(msg) { can.user.toastSuccess(can) }) },
	request: function(event, can, button) { can.runAction(can.onexport.request(event, can), button, [], function(msg) {
		delete(can.ui.response.content._cache), delete(can.ui.response.content._cache_key)
		can.db.msg = msg, can.page.SelectOne(can, can.ui.response.action, "").click()
		if (msg.IsErr()) {
			can.user.toastFailure(can, msg.Result())
		} else {
			can.user.toastSuccess(can)
		}
	}) },
	profile: function(event, can, button) { can.onmotion.toggle(can, can.ui.profile), can.onimport.layout(can) },
	display: function(event, can, button) { can.onmotion.toggle(can, can.ui.display), can.onimport.layout(can) },
	enter: function(event, can) { can.onaction.request(event, can, REQUEST) },
})
Volcanos(chat.ONEXPORT, {
	request: function(event, can) { var msg = can.request(event, can.db.current, can.Option())
		can.page.Select(can, can.ui.action, "select[name=method]", function(target) { msg.Option(METHOD, target.value) })
		can.page.Select(can, can.ui.action, "input[name=url]", function(target) { msg.Option("url", target.value) })
		var _select = can.page.SelectOne(can, can.ui.request.action, html.DIV_ITEM_SELECT)
		can.page.Select(can, can.ui.request.action, "", function(target) { var args = {}; target.click()
			can.page.Select(can, can.ui.request.content, html.TR, function(tr, index) { if (index == 0) { return }
				var input = can.page.Select(can, tr, html.INPUT); input[0].value && (args[input[0].value] = input[1].value)
			}), msg.Option(target._item.name, JSON.stringify(args))
		}), _select.click(); return msg
	},
	config: function(event, can) { var msg = can.onexport.request(event, can); return can.base.Obj(msg.Option(ctx.CONFIG))||{} },
	_args: function(can, key) { 
		var args = [], config = can.onexport.config(event, can), value = config[key]; if (!value) { return }
		if (can.base.beginWith(value, nfs.PS, web.HTTP)) {
			args = [web.SPIDE, value]
		} else if (can.base.beginWith(value, nfs.SRC, nfs.USR)) {
			args = can.misc.SplitPath(can, value)
		} else {
			args = [can.misc.Template(can, ctx.CONFIG, key)+nfs.PS, value]
		}
		return args
	},
	tabs: function(can) { can.misc.sessionStorage(can, [can.ConfIndex(), html.TABS], can.page.Select(can, can._action, html.DIV_TABS, function(target) { return target._item.hash })) },
})
})()
