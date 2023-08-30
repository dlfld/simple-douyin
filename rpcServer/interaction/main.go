package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/jaeger"
	interaction "github.com/douyin/kitex_gen/interaction/interactionservice"
)

func main() {
	tracerSuite, closer := jaeger.InitJaegerServer("kitex-server-interaction")
	defer closer.Close()
	addr, err := net.ResolveTCPAddr("tcp", conf.InteractionService.Addr)
	if err != nil {
		log.Printf("addr获取失败：%+v\n", err)
	}

	svr := interaction.NewServer(new(InteractionServiceImpl), server.WithServiceAddr(addr), server.WithSuite(tracerSuite))
	InitDao()

	err = svr.Run()
	if err != nil {
		log.Printf("rpc服务启动失败：%+v\n", err)
	}

}
