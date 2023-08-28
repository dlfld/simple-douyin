// Package video /*
package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/kafkaLog/productor"
	video "github.com/douyin/kitex_gen/video/videoservice"
	"log"
	"net"
)

// LogCollector 日志收集器
var LogCollector *productor.LogCollector

func init() {
	var err error
	//初始化日志收集器
	if LogCollector, err = productor.NewLogCollector(conf.MessageService.Name); err != nil {
		panic(err)
	}
}
func main() {
	addr, err := net.ResolveTCPAddr("tcp", conf.VideoService.Addr)
	svr := video.NewServer(new(VideoServiceImpl), server.WithServiceAddr(addr))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
