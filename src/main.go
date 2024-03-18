package main

import (
	"shylinux.com/x/ice"
	_ "shylinux.com/x/ice/devops"

	_ "shylinux.com/x/nginx-story/src/client"
	_ "shylinux.com/x/nginx-story/src/consul"
	_ "shylinux.com/x/nginx-story/src/server"
	_ "shylinux.com/x/nginx-story/src/studio"
)

func main() { print(ice.Run()) }
