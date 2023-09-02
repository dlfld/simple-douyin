/*
*

	@author:孟令亚
	@date:2023/8/9
	@node:

*
*/
package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/etcd"
	"github.com/douyin/common/jaeger"
	"github.com/douyin/common/kafkaLog/productor"
	"github.com/douyin/kitex_gen/user/userservice"
)

var logCollector *productor.LogCollector

func init() {
	var err error
	if logCollector, err = productor.NewLogCollector(conf.UserService.Name); err != nil {
		panic(err)
	}
}

func main() {
	tracerSuite, closer := jaeger.InitJaegerServer("kitex-server-user")
	defer closer.Close()
	addr, err := net.ResolveTCPAddr("tcp", conf.UserService.Addr)
	svr := userservice.NewServer(new(UserServiceImpl), server.WithServiceAddr(addr), server.WithSuite(tracerSuite))
	if err != nil {
		log.Println(err.Error())
	}
	err = svr.Run()
	etcd.RegisterService(conf.UserService.Name, conf.UserService.Addr)
	if err != nil {
		log.Println(err.Error())
	}

}
