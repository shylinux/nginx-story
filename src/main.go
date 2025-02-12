package main

import (
	"shylinux.com/x/ice"
	_ "shylinux.com/x/nginx-story/src/minio"
	_ "shylinux.com/x/nginx-story/src/server"
	_ "shylinux.com/x/nginx-story/src/studio"
)

func main() { print(ice.Run()) }

func init() {
	ice.Info.NodeIcon = "src/server/nginx.png"
	ice.Info.NodeMain = "web.code.nginx.configs"
}
