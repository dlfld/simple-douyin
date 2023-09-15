package main

import (
	"context"
	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/bloom"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/etcd"
	"github.com/douyin/common/kafkaLog/productor"
	"github.com/douyin/common/otel"
	"github.com/douyin/kitex_gen/relation/relationservice"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"

	"log"
	"net"
)

var logCollector *productor.LogCollector

func init() {
	logCollector = &productor.LogCollector{ServiceName: conf.RelationService.Name}
}
func main() {
	p := otel.NewOtelProvider("relation")
	defer p.Shutdown(context.Background())
	var err error
	if logCollector, err = productor.NewLogCollector(conf.RelationService.Name); err != nil {
		panic(err)
	}

	addr, err := net.ResolveTCPAddr("tcp", conf.RelationService.Addr)
	svr := relationservice.NewServer(new(RelationServiceImpl), server.WithServiceAddr(addr), server.WithSuite(tracing.NewServerSuite()))
	if err != nil {
		log.Println(err.Error())
	}
	etcd.RegisterService(conf.RelationService.Name, conf.RelationService.Addr)
	bf = bloom.NewBloom()
	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}

}
