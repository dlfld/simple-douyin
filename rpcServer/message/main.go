package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	message "github.com/douyin/kitex_gen/message/messageservice"
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", conf.MessageService.Addr)
	svr := message.NewServer(new(MessageServiceImpl), server.WithServiceAddr(addr))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
