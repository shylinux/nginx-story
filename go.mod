module shylinux.com/x/nginx-story

go 1.13

replace (
	shylinux.com/x/ice => ./usr/release
	shylinux.com/x/icebergs => ./usr/icebergs
	shylinux.com/x/toolkits => ./usr/toolkits
)

require (
	shylinux.com/x/ice v1.5.0
	shylinux.com/x/icebergs v1.9.0
	shylinux.com/x/toolkits v1.0.4
)

require github.com/hashicorp/consul/api v1.24.0
