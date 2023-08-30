// Package video /*
package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/jaeger"
	"github.com/douyin/common/kafkaLog/productor"
	"github.com/douyin/common/etcd"
	video "github.com/douyin/kitex_gen/video/videoservice"
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
	etcd.RegisterService(conf.VideoService.Name, conf.VideoService.Addr)
}
