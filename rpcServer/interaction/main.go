package main

import (
	"context"
	"log"
	"net"

	"github.com/douyin/common/etcd"
	"github.com/douyin/common/otel"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"

	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/kafkaLog/productor"
	interaction "github.com/douyin/kitex_gen/interaction/interactionservice"
)

func main() {
	p := otel.NewOtelProvider("interaction")
	defer p.Shutdown(context.Background())
	addr, err := net.ResolveTCPAddr("tcp", conf.InteractionService.Addr)

	if err != nil {
		log.Printf("addr获取失败：%+v\n", err)
	}

	svr := interaction.NewServer(new(InteractionServiceImpl), server.WithServiceAddr(addr), server.WithSuite(tracing.NewServerSuite()))
	InitDao()
	etcd.RegisterService(conf.InteractionService.Name, conf.InteractionService.Addr)
	if logCollector, err = productor.NewLogCollector(conf.InteractionService.Name); err != nil {
		panic(err)
	}
	logCollector.Info("Interaction 服务启动")
	err = svr.Run()
	if err != nil {
		panic(err)
	}

}
