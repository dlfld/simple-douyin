package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/etcd"
	interaction "github.com/douyin/kitex_gen/interaction/interactionservice"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", conf.InteractionService.Addr)

	if err != nil {
		log.Printf("addr获取失败：%+v\n", err)
	}

	svr := interaction.NewServer(new(InteractionServiceImpl), server.WithServiceAddr(addr))
	InitDao()

	err = svr.Run()
	if err != nil {
		panic(err)
	}
	etcd.RegisterService(conf.InteractionService.Name, conf.InteractionService.Addr)

}
