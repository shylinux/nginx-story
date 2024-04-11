package main

import (
	"shylinux.com/x/ice"
	_ "shylinux.com/x/ice/devops"
	_ "shylinux.com/x/nginx-story/src/server"
	_ "shylinux.com/x/nginx-story/src/studio"
)

func init() { ice.Info.NodeIcon = "src/server/nginx.png" }

func main() { print(ice.Run()) }
