(function() {
const REQUEST = "request", RESPONSE = "response"
const METHOD = "method", PARAMS = "params", STATUS = "status", HEADER = "header", COOKIE = "cookie"
const PLUGIN_STORY_EDITOR = "/plugin/story/editor.js"
const PLUGIN_STORY_MONACO = "/plugin/story/monaco.js"
const PLUGIN_STORY_JSON = "/plugin/story/json.js"
Volcanos(chat.ONIMPORT, {
	_init: function(can, msg) { can.ui = can.onappend.layout(can), can.onimport._project(can, msg) },
	_project: function(can, msg) {
		msg.Table(function(value) { value.type = value.method
			value.nick = can.page.Format(html.SPAN, value.method == http.DELETE? "DEL": value.method, [METHOD, value.method])+lex.SP+can.page.Format(html.SPAN, value.name, mdb.NAME)
			value._select = can.db.hash[0] == value.hash
			can.onimport.item(can, value, function(event, value, show, target) {
				show == undefined && can.onimport.tabsCache(can, can.request(event), value.hash, value, target, function() {
					can.onimport._content(can, value)
				})
			})
		})
	},
	_content: function(can, value) { value._msg = can.request()
		value._action = can.page.Append(can, can.ui.content, [
			{view: html.ACTION, _init: function(target) { can.onappend._action(can, [
				{type: html.SELECT, name: METHOD, value: value.method, values: [http.GET, http.PUT, http.POST, http.DELETE]},
				{type: html.TEXT, name: web.URL, value: value.url, action: "key"},
				{type: html.BUTTON, value: REQUEST, style: html.NOTICE},
			], target) }},
		])._target
		value._request = can.onimport._part(can, REQUEST, can.core.List([PARAMS, HEADER, COOKIE, aaa.AUTH, ctx.CONFIG], function(key) {
			return {name: key, show: function(event, target) {
				var table = can.onimport._table(can, target, value, key)
				var action = can.page.Append(can, target, [html.ACTION])._target
				can.onappend._action(can, ["addline"], action, {
					"addline": function() {
						can.page.Select(can, table, html.TBODY, function(target) {
							can.page.Append(can, target, [{type: html.TR, list: [table._add("", mdb.NAME), table._add("", mdb.VALUE)]}])
						})
					},
				})
			}}
		}), can.ui.content, [{type: html.BUTTON, value: nfs.SAVE}])
		value._response = can.onimport._part(can, RESPONSE, [{name: mdb.DATA, show: function(event, target) {
			can.onimport._plugin(can, target)
		}}].concat(can.core.List([STATUS, HEADER, COOKIE], function(key) {
			return {name: key, show: function(event, target) { can.onimport._response(can, target, can.db.value._msg, key) }}
		})), can.ui.content, [{type: html.BUTTON, value: REQUEST, style: html.NOTICE}])
	},
	_part: function(can, name, list, target, action) {
		var ui = can.page.Append(can, target, [
			{view: name, list: [
				{view: html.TITLE, list: [{text: can.base.capital(name)}], _init: function(target) {
					can.onappend._action(can, action, target, null, true)
				}},
				{view: html.ACTION, list: can.core.List(list, function(item) {
					return {view: [[html.ITEM, item.name], "", item.name], onclick: function(event) {
						can.onmotion.select(can, ui.action, "", event.currentTarget)
						if (can.onmotion.cache(can, function() { return item.name}, ui.content)) { return }
						item.show(event, ui.content)
					}, _init: function(target) { target._item = item }}
				})},
				{view: html.CONTENT},
			]},
		]); can.page.SelectOne(can, ui.action, "", function(target) { target.click() }); return ui
	},
	_table: function(can, target, value, _key, keys) {
		var msg = can.request(); keys = keys || [mdb.NAME, mdb.VALUE]
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
		}, list: []} } var table = can.onappend.table(can, msg, add, target); table._add = add; return table
	},
	_response: function(can, target, msg, type) { var _msg = can.request()
		msg.Table(function(value) { if (value.type == type) { _msg.Push(mdb.NAME, value.name), _msg.Push(mdb.VALUE, value.value) } })
		can.onappend.style(can, type, can.onappend.table(can, _msg, null, target))
	},
	_plugin: function(can, target) {
		can.db.value._msg.Table(function(value) { if (value.name == http.ContentType) {
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
					index: config.index, args: config.args, msg: can.db.value._msg, style: config.style, display: config.display, height: can.ConfHeight()/2-1,
				}, function(sub) { can.ui.content._plugin = sub
					can.onmotion.hidden(can, sub._legend), can.onmotion.hidden(can, sub._option)
				}, target)
				return
			}
			switch (can.core.Split(value.value, ";")[0]) {
				case http.ApplicationJSON: can.onappend._output(can, can.db.value._msg, PLUGIN_STORY_JSON, null, target, null, false); break
				default: can.onappend.table(can, can.db.value._msg, null, target), can.onappend.board(can, can.db.value._msg, target)
			}
		} })
	},
})
Volcanos(chat.ONACTION, {
	save: function(event, can, button) { var msg = can.onexport.request(event, can)
		can.runAction(msg, mdb.MODIFY, msg.OptionSimple(mdb.HASH, METHOD, "url", "params", "header", "cookie"), function(msg) { can.user.toastSuccess(can) }) },
	request: function(event, can, button) { can.runAction(can.onexport.request(event, can), button, [], function(msg) {
		delete(can.db.value._response.content._cache), delete(can.db.value._response.content._cache_key)
		can.db.value._msg = msg, can.page.SelectOne(can, can.db.value._response.action, "").click()
	}) },
})
Volcanos(chat.ONEXPORT, {
	request: function(event, can) { var msg = can.request(event, can.db.value, can.Option())
		can.page.Select(can, can.db.value._action, "select[name=method]", function(target) { msg.Option(METHOD, target.value) })
		can.page.Select(can, can.db.value._action, "input[name=url]", function(target) { msg.Option("url", target.value) })
		var _select = can.page.SelectOne(can, can.db.value._request.action, html.DIV_ITEM_SELECT)
		can.page.Select(can, can.db.value._request.action, "", function(target) { var args = {}; target.click()
			can.page.Select(can, can.db.value._request.content, html.TR, function(tr, index) { if (index == 0) { return }
				var input = can.page.Select(can, tr, html.INPUT); input[0].value && (args[input[0].value] = input[1].value)
			}), msg.Option(target._item.name, JSON.stringify(args))
		}), _select.click(); return msg
	},
	config: function(event, can) { var msg = can.onexport.request(event, can); return can.base.Obj(msg.Option(ctx.CONFIG))||{} },
})
})()
