/*
*

	@author:戴林峰
	@date:2023/7/29
	@node:

*
*/
package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/etcd"
	video "github.com/douyin/kitex_gen/video/videoservice"
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", conf.VideoService.Addr)
	svr := video.NewServer(new(VideoServiceImpl), server.WithServiceAddr(addr))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
	etcd.RegisterService(conf.VideoService.Name, conf.VideoService.Addr)
}
