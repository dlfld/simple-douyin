package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	interaction "github.com/douyin/kitex_gen/interaction/interactionservice"
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", conf.InteractionService.Addr)
	if  err != nil {
		log.Printf("addr获取失败：%+v\n",err)
	}
	svr := interaction.NewServer(new(InteractionServiceImpl), server.WithServiceAddr(addr))
	err = svr.Run()
	if err != nil {
		log.Printf("rpc服务启动失败：%+v\n", err)
	}
}
