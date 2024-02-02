Volcanos(chat.ONSYNTAX, {
	conf: {
		prefix: {"#": code.COMMENT},
		keyword: {
		"server": code.KEYWORD,
		"location": code.KEYWORD,
		"upstream": code.KEYWORD,
		"include": code.KEYWORD,

		"listen": code.FUNCTION,
		"server_name": code.FUNCTION,
		"proxy_pass": code.FUNCTION,
	},func: function(can, push, text, indent, opts) {
	}},
})
