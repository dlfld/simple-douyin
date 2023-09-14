package relationRpc

import (
	"github.com/douyin/common/etcd"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/jaeger"
	"github.com/douyin/kitex_gen/relation/relationservice"
)

var cli relationservice.Client
var err error
var once sync.Once

func init() {
	tracerSuite, _ := jaeger.InitJaeger("kitex-client-relation")
	cli, err = relationservice.NewClient(conf.RelationService.Name, client.WithHostPorts(conf.RelationService.Addr), client.WithSuite(tracerSuite))
	if err != nil {
		panic(err)
	}
}

func NewRpcRelationClient() (relationservice.Client, error) {
	once.Do(func() {
		addr := etcd.DiscoverService(conf.RelationService.Name)
		cli, err = relationservice.NewClient(conf.RelationService.Name, client.WithHostPorts(addr...), client.WithSuite(tracing.NewClientSuite()))
		if err != nil {
			panic(err)
		}
	})
	return cli, nil
}
