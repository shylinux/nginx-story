package main

import (
	"shylinux.com/x/ice"
	_ "shylinux.com/x/nginx-story/src/client"
	_ "shylinux.com/x/nginx-story/src/consul"
	_ "shylinux.com/x/nginx-story/src/server"
)

func main() { print(ice.Run()) }
