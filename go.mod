module shylinux.com/x/nginx-story

go 1.13

replace (
	shylinux.com/x/ice => ./usr/release
	shylinux.com/x/icebergs => ./usr/icebergs
	shylinux.com/x/toolkits => ./usr/toolkits
)

require (
	shylinux.com/x/ice v1.3.11
	shylinux.com/x/icebergs v1.5.18
	shylinux.com/x/toolkits v0.7.9
)

require (
	github.com/araddon/gou v0.0.0-20211019181548-e7d08105776c // indirect
	github.com/hashicorp/consul/api v1.21.0
	github.com/lytics/confl v0.0.0-20200313154245-08c6aed5f53f // indirect
	github.com/tufanbarisyildirim/gonginx v0.0.0-20230627120331-964b6ae8380e // indirect
	github.com/webantic/nginx-config-parser v1.1.0 // indirect
)
