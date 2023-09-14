// Package video /*
package main

import (
	"context"
	"github.com/douyin/common/bloom"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/etcd"
	"github.com/douyin/common/kafkaLog/productor"
	video "github.com/douyin/kitex_gen/video/videoservice"
)

// LogCollector 日志收集器
var LogCollector *productor.LogCollector
var bf *bloom.Filter

func init() {
	var err error
	//初始化日志收集器
	if LogCollector, err = productor.NewLogCollector(conf.MessageService.Name); err != nil {
		panic(err)
	}
	bf = bloom.NewBloom()
}
func main() {
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName("video"),
		provider.WithExportEndpoint("host.docker.internal:4317"),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	addr, err := net.ResolveTCPAddr("tcp", conf.VideoService.Addr)
	if err != nil {
		log.Println(err.Error())
	}
	svr := video.NewServer(new(VideoServiceImpl), server.WithServiceAddr(addr), server.WithSuite(tracing.NewServerSuite()))
	etcd.RegisterService(conf.VideoService.Name, conf.VideoService.Addr)
	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}

}
