package main

import (
	ice "shylinux.com/x/icebergs"
	_ "shylinux.com/x/icebergs/base"
	_ "shylinux.com/x/icebergs/core"
	_ "shylinux.com/x/icebergs/misc"

	_ "shylinux.com/x/nginx-story/src/server"
)

func main() { print(ice.Run()) }
