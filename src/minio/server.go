package minio

import "shylinux.com/x/ice"

type server struct {
	ice.Code
	// linux string `data:"https://dl.min.io/server/minio/release/linux-amd64/minio"`
	linux string `data:"https://2024-jingganjiaoyu.shylinux.com/publish/minio"`
}

func (s server) List(m *ice.Message, arg ...string) {
	s.Code.List(m, "", arg...)
}

func init() { ice.Cmd("web.chat.minio.server", server{}) }
