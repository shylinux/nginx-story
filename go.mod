module shylinux.com/x/nginx-story

go 1.13

replace (
	shylinux.com/x/ice => ./usr/release
	shylinux.com/x/icebergs => ./usr/icebergs
	shylinux.com/x/toolkits => ./usr/toolkits
)

require (
	github.com/pkg/sftp v1.13.6 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	shylinux.com/x/ice v1.5.4
	shylinux.com/x/icebergs v1.9.4
	shylinux.com/x/toolkits v1.0.4
)
