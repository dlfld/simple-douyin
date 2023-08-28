/*
*

	@author:戴林峰
	@date:2023/7/29
	@node:

*
*/
package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/jaeger"
	video "github.com/douyin/kitex_gen/video/videoservice"
)

func main() {
	tracerSuite, closer := jaeger.InitJaegerServer("kitex-server-video")
	defer closer.Close()
	addr, err := net.ResolveTCPAddr("tcp", conf.VideoService.Addr)
	if err != nil {
		log.Println(err.Error())
	}
	svr := video.NewServer(new(VideoServiceImpl), server.WithServiceAddr(addr), server.WithSuite(tracerSuite))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
