package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	interaction "github.com/douyin/kitex_gen/interaction/interactionservice"
)

func main() {
	err := InitDao()
	if err != nil {
		log.Printf("新建数据库连接失败：%+v\n", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", conf.InteractionService.Addr)
	if err != nil {
		log.Printf("addr获取失败：%+v\n", err)
	}

	svr := interaction.NewServer(new(InteractionServiceImpl), server.WithServiceAddr(addr))

	err = svr.Run()
	if err != nil {
		log.Printf("rpc服务启动失败：%+v\n", err)
	}
}
