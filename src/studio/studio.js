(function() { const REQUEST = "request", RESPONSE = "response"
const METHOD = "method", PARAMS = "params", STATUS = "status", HEADER = "header", COOKIE = "cookie"
Volcanos(chat.ONIMPORT, {
	_init: function(can, msg, cb) { can.onmotion.clear(can), can.ui = can.onappend.layout(can), cb && cb(msg)
		can.db._list = {}, can.db._hash = can.misc.SearchHash(can), can.onimport._project(can, msg)
	},
	_project: function(can, msg) {
		var _select; msg.Table(function(value) {
			if (value.method == web.DELETE) {
				value.nick = can.page.Format(html.SPAN, "DEL", METHOD+lex.SP+"DEL")+lex.SP+value.name
			} else {
				value.nick = can.page.Format(html.SPAN, value.method, METHOD+lex.SP+value.method)+lex.SP+value.name
			}
			var _target = can.onimport.item(can, value, function(event) { if (value._tabs) { return value._tabs.click() }
				value._tabs = can.onimport.tabs(can, [value], function(event) { can.db.current = value, can.misc.SearchHash(can, value.hash)
					if (can.onmotion.cache(can, function(save, load) {
						save({msg: can.db.msg, request: can.ui.request, response: can.ui.response})
						return load(value.hash, function(bak) { can.db.msg = bak.msg, can.ui.request = bak.request, can.ui.response = bak.response })
					}, can.ui.content)) { return can.onimport.layout(can), _target.click() }
					can.onimport._content(can, value), can.onimport.layout(can)
				}, function(event) { delete(can.ui.content._cache[value.hash])
					can.onmotion.delay(can, function() { can.onexport.tabs(can) })
				}), can.onexport.tabs(can)
			}, null, can.ui.project); can.db._list[value.hash] = _target, (!_select || value.hash == can.db._hash[0]) && (_select = _target)
		})
		can.core.List(can.misc.sessionStorage(can, [can.ConfIndex(), html.TABS]), function(hash) {
			var tabs = can.db._list[hash]; tabs && tabs.click()
		}), _select && _select.click()
	},
	_content: function(can, value) { can.db.msg = can.request()
		can.page.Append(can, can.ui.content, [
			{view: html.ACTION, _init: function(target) { can.onappend._action(can, [
				{type: html.SELECT, name: METHOD, value: value.method, values: [web.GET, web.PUT, web.POST, web.DELETE]},
				{type: html.TEXT, name: web.URL, value: value.url, action: "key"},
				{type: html.BUTTON, value: REQUEST, style: html.NOTICE},
				{type: html.BUTTON, value: nfs.SAVE},
			], target) }},
		])
		value.description && (can.onimport._part(can, "description", [], can.ui.content).content.innerHTML = value.description)
		can.ui.request = can.onimport._part(can, REQUEST, can.core.List([PARAMS, HEADER, COOKIE, aaa.AUTH, ctx.CONFIG], function(key) {
			return {name: key, show: function(event, target) { can.onimport._table(can, target, value, key) }}
		}), can.ui.content, [{type: html.BUTTON, value: nfs.SAVE}])
		can.ui.response = can.onimport._part(can, RESPONSE, [{name: mdb.DATA, show: function(event, target) {
			can.onimport._plugin(can, target)
		}}].concat(can.core.List([STATUS, HEADER, COOKIE], function(key) {
			return {name: key, show: function(event, target) { can.onimport._response(can, target, can.db.msg, key) }}
		}), [{name: ctx.DISPLAY, show: function(event, target) {
			var msg = can.onexport.request(event, can), display = can.core.Value(can.base.Obj(msg.Option(ctx.CONFIG)), ctx.DISPLAY)
			return can.onappend.plugin(can, {index: web.CODE_INNER, args: [display||"/volcanos/plugin/story/json.js", lex.SP], style: html.OUTPUT}, function(sub) {}, target)
		}}]), can.ui.content, [{type: html.BUTTON, value: REQUEST, style: html.NOTICE}])
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
	_table: function(can, target, value, _key, keys) { var msg = can.request(); keys = keys || [mdb.NAME, mdb.VALUE, "description"]
		can.core.Item(can.base.Obj(value[_key]), function(key, value) { msg.Push(mdb.NAME, key), msg.Push(mdb.VALUE, value) })
		can.core.List(keys, function(key) { msg.Push(key, "") })
		function add(value, key) { return {type: html.TD, _init: function(target) {
			can.onappend._action(can, [{type: html.TEXT, name: key, value: value, _init: function(target) {
				can.onappend.figure(can, {name: key, run: function(event, cmds, cb) { var msg = can.request(event, {action: _key})
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
		var msg = can.onexport.request({}, can), display = can.core.Value(can.base.Obj(msg.Option(ctx.CONFIG)), ctx.DISPLAY)
		can.db.msg.Table(function(value) { if (value.name == http.ContentType) {
			if (display) {
				return can.onappend.plugin(can, {msg: can.db.msg, display: display}, function(sub) {
					can.onmotion.hidden(can, sub._legend), can.onmotion.hidden(can, sub._option)
				}, target)
			}
			switch (can.core.Split(value.value, ";")[0]) {
				case mime.ApplicationJSON: can.onappend._output(can, can.db.msg, "/plugin/story/json.js", null, target, null, false); break
				default: can.onappend.table(can, can.db.msg, null, target), can.onappend.board(can, can.db.msg, target)
			}
		} }) 
	},
}, [""])
Volcanos(chat.ONACTION, {
	save: function(event, can, button) { can.runAction(can.onexport.request(event, can), button) },
	request: function(event, can, button) { can.runAction(can.onexport.request(event, can), button, [], function(msg) {
		delete(can.ui.response.content._cache), delete(can.ui.response.content._cache_key)
		can.db.msg = msg, can.page.SelectOne(can, can.ui.response.action, "").click()
		can.user.toastSuccess(can)
	}) },
	enter: function(event, can) { can.onaction.request(event, can, REQUEST) },
})
Volcanos(chat.ONEXPORT, {
	request: function(event, can) { var msg = can.request(event, can.db.current, can.Option())
		var _select = can.page.SelectOne(can, can.ui.request.action, html.DIV_ITEM_SELECT)
		can.page.Select(can, can.ui.request.action, "", function(target) { var args = {}; target.click()
			can.page.Select(can, can.ui.request.content, html.TR, function(tr, index) { if (index == 0) { return }
				var input = can.page.Select(can, tr, html.INPUT); input[0].value && (args[input[0].value] = input[1].value)
			}), msg.Option(target._item.name, JSON.stringify(args))
		}), _select.click(); return msg
	},
	tabs: function(can) {
		can.misc.sessionStorage(can, [can.ConfIndex(), html.TABS], can.page.Select(can, can._action, html.DIV_TABS, function(target) { return target._item.hash }))
	},
})
})()
