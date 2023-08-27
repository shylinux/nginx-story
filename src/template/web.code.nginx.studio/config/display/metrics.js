Volcanos(chat.ONIMPORT, {
	_init: function(can, msg, cb) {
		var _msg = can.request()
		can.core.List(msg.Result().split("\n"), function(item) {
			item = item.trim()
			if (item.indexOf("#") == 0) { return }
			var ls = can.core.Split(item, " {=,}", )
			_msg.Push(mdb.NAME, ls[0])
			_msg.Push(mdb.VALUE, ls.pop())
		})
		can.onappend.table(can, _msg)
		can.onappend.board(can, msg)
		_msg.StatusTimeCount()
		cb && cb(_msg)
	},
})

