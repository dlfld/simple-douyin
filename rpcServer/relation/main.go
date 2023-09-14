package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/etcd"
	"github.com/douyin/common/jaeger"
	"github.com/douyin/common/kafkaLog/productor"
	"github.com/douyin/kitex_gen/relation/relationservice"

	"log"
	"net"
)

var logCollector *productor.LogCollector

func init() {
	logCollector = &productor.LogCollector{ServiceName: conf.RelationService.Name}
}
func main() {
	tracerSuite, closer := jaeger.InitJaegerServer("kitex-server-relation")
	var err error
	if logCollector, err = productor.NewLogCollector(conf.RelationService.Name); err != nil {
		panic(err)
	}
	defer closer.Close()

	addr, err := net.ResolveTCPAddr("tcp", conf.RelationService.Addr)
	svr := relationservice.NewServer(new(RelationServiceImpl), server.WithServiceAddr(addr), server.WithSuite(tracerSuite))
	if err != nil {
		log.Println(err.Error())
	}
	etcd.RegisterService(conf.RelationService.Name, conf.RelationService.Addr)
	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}

}
