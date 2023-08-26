(function() { const REQUEST = "request", RESPONSE = "response"
const METHOD = "method", PARAMS = "params", STATUS = "status", HEADER = "header", COOKIE = "cookie"
Volcanos(chat.ONIMPORT, {
	_init: function(can, msg, cb) {
		can.onmotion.clear(can), can.ui = can.onappend.layout(can), can.onmotion.hidden(can, can.ui.profile), can.onmotion.hidden(can, can.ui.display)
		can.db._list = {}, can.db._hash = can.misc.SearchHash(can), cb && cb(msg)
		var _select; msg.Table(function(value) {
			if (value.method == web.DELETE) {
				value.nick = `<span class="method DEL">DEL </span>`+value.name
			} else {
				value.nick = `<span class="method ${value.method}">${value.method} </span>`+value.name
			}
			var _target = can.onimport.item(can, value, function(event) { if (value._tabs) { return value._tabs.click() }
				value._tabs = can.onimport.tabs(can, [value], function(event) { can.db.current = value, can.misc.SearchHash(can, value.hash)
					if (can.onmotion.cache(can, function() { return value.hash }, can.ui.content)) { return can.onimport.layout(can) }
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
	_content: function(can, value) { var _msg = can.request()
		can.page.Append(can, can.ui.content, [
			{view: html.ACTION, _init: function(target) { can.onappend._action(can, [
				{type: html.SELECT, name: METHOD, value: value.method, values: [web.GET, web.PUT, web.POST, web.DELETE]},
				{type: html.INPUT, name: "url", value: value.url},
				{type: html.BUTTON, value: REQUEST, onclick: function(event) {
					can.run(can.onexport.request(can, request, value), [ctx.ACTION, REQUEST], function(msg) {
						delete(response.content._cache), delete(response.content._cache_key)
						_msg = msg, can.page.SelectOne(can, response.action, "").click()
						can.user.toastSuccess(can)
					})
				}},
				{type: html.BUTTON, value: nfs.SAVE, onclick: function(event) { can.runAction(can.onexport.request(can, request, value), nfs.SAVE) }},
				{type: html.BUTTON, value: mdb.DELETE, onclick: function(event) { can.runAction(can.onexport.request(can, request, value), mdb.DELETE) }}, 
			], target) }},
		])
		var request = can.onimport._part(can, REQUEST, can.core.List([PARAMS, HEADER, COOKIE, "auth"], function(key) {
			return {name: key, show: function(event, target) { can.onimport._table(can, target, value, key) }}
		}), can.ui.content)
		var response = can.onimport._part(can, RESPONSE, [
			{name: "data", show: function(event, target) { can.onappend._output(can, _msg, "/plugin/story/json.js", target, false) }},
		].concat(can.core.List([STATUS, HEADER, COOKIE], function(key) {
			return {name: key, show: function(event, target) { can.onimport._response(can, target, _msg, key) }}
		})), can.ui.content)
	},
	_part: function(can, name, list, target) {
		var ui = can.page.Append(can, target, [
			{view: name, list: [{view: [html.TITLE, "", name]},
				{view: html.ACTION, list: can.core.List(list, function(item) {
					return {view: [[html.ITEM, item.name], "", item.name], onclick: function(event) {
						can.onmotion.select(can, ui.action, "", event.currentTarget)
						if (can.onmotion.cache(can, function() { return item.name}, ui.content)) { return }
						item.show(event, ui.content)
					}, _init: function(target) { target._item = item }}
				})},
				{view: html.CONTENT},
			]},
		]); can.page.SelectOne(can, ui.action, "").click()
		return ui
	},
	_table: function(can, target, value, _key, keys) { var msg = can.request(); keys = keys || [mdb.NAME, mdb.VALUE, "description"]
		can.core.Item(can.base.Obj(value[_key]), function(key, value) { msg.Push(mdb.NAME, key), msg.Push(mdb.VALUE, value) })
		can.core.List(keys, function(key) { msg.Push(key, "") })
		function add(value, key) { return {type: html.TD, list: [{type: html.INPUT, value: value, _init: function(target) {
			can.onappend.figure(can, {run: function(event, cmds, cb) { var msg = can.request(event, {action: _key})
				can.page.Select(can, target.parentNode.parentNode, html.INPUT, function(target, index) { msg.Option(keys[index], target.value) })
				can.run(event, [ctx.ACTION, mdb.INPUTS, key], cb)
			}, _enter: function() {
				can.page.Append(can, table, [{type: html.TR, list: can.core.List(keys, function(key) { return add("", key) }) }])
			}}, target, function(sub, value) {})
		}}]} }
		var table = can.onappend.table(can, msg, add, target); return table
	},
	_response: function(can, target, msg, type) { var _msg = can.request()
		msg.Table(function(value) { if (value.type == type) { _msg.Push(mdb.NAME, value.name), _msg.Push(mdb.VALUE, value.value) } })
		can.onappend.table(can, _msg, null, target)
	},
	layout: function(can) { can.ui.layout(can.ConfHeight(), can.ConfWidth()) },
}, [""])
Volcanos(chat.ONEXPORT, {
	request: function(can, request, value) { var msg = can.request(event, value, can.Option())
		can.page.Select(can, request.action, "", function(target) { var args = {}; target.click()
			can.page.Select(can, request.content, html.TR, function(tr, index) { if (index == 0) { return }
				var input = can.page.Select(can, tr, html.INPUT); input[0].value && (args[input[0].value] = input[1].value)
			}), msg.Option(target._item.name, JSON.stringify(args))
		}), can.page.SelectOne(can, request.action).click(); return msg
	},
	tabs: function(can) {
		can.misc.sessionStorage(can, [can.ConfIndex(), html.TABS], can.page.Select(can, can._action, html.DIV_TABS, function(target) { return target._item.hash }))
	},
})
})()
