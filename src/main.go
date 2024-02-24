package main

import (
	"shylinux.com/x/ice"
	_ "shylinux.com/x/icebergs/base/aaa/portal"
	_ "shylinux.com/x/icebergs/core/chat/oauth"
	_ "shylinux.com/x/icebergs/misc/java"
	_ "shylinux.com/x/icebergs/misc/node"
	_ "shylinux.com/x/icebergs/misc/wx"

	_ "shylinux.com/x/nginx-story/src/client"
	_ "shylinux.com/x/nginx-story/src/consul"
	_ "shylinux.com/x/nginx-story/src/server"
	_ "shylinux.com/x/nginx-story/src/studio"
)

func main() { print(ice.Run()) }
