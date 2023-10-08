package main

import (
	"fmt"
	"github.com/cloudwego/kitex/server"
	"github.com/douyin/kitex_gen/gpt/chatgptservice"
	"log"
	"net"
)

var GptService = "gpt-service"
var GptAddr = "0.0.0.0:8888"

func main() {
	addr, err := net.ResolveTCPAddr("tcp", GptAddr)
	if err != nil {
		log.Println(err.Error())
	}
	svr := chatgptservice.NewServer(new(ChatgptServiceImpl), server.WithServiceAddr(addr))
	fmt.Println("start gpt service")
	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
