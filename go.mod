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

require github.com/hashicorp/consul/api v1.21.0
