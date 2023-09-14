package interactionRpc

import (
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/etcd"
	"github.com/douyin/kitex_gen/interaction/interactionservice"
)

var once sync.Once
var cli interactionservice.Client
var err error

func NewRpcInteractionClient() (interactionservice.Client, error) {
	once.Do(func() {
		addr := etcd.DiscoverService(conf.InteractionService.Name)
		cli, err = interactionservice.NewClient(conf.InteractionService.Name, client.WithHostPorts(addr...), client.WithSuite(tracing.NewClientSuite()))
		if err != nil {
			panic(err)
		}
	})
	return cli, err
}
