package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	"github.com/douyin/kitex_gen/relation/relationservice"

	"log"
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", conf.RelationService.Addr)
	svr := relationservice.NewServer(new(RelationServiceImpl), server.WithServiceAddr(addr))
	if err != nil {
		log.Println(err.Error())
	}
	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
