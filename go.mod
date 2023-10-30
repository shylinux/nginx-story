module shylinux.com/x/nginx-story

go 1.13

replace (
	shylinux.com/x/ice => ./usr/release
	shylinux.com/x/icebergs => ./usr/icebergs
	shylinux.com/x/toolkits => ./usr/toolkits
)

require (
	github.com/hashicorp/consul/api v1.24.0
	shylinux.com/x/ice v1.3.16
	shylinux.com/x/icebergs v1.6.4
	shylinux.com/x/toolkits v0.7.13
)
