package interactionRpc

import (
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/jaeger"
	"github.com/douyin/kitex_gen/interaction/interactionservice"
)

var once sync.Once
var cli interactionservice.Client
var err error

func NewRpcInteractionClient() (interactionservice.Client, error) {
	once.Do(func() {
		tracerSuite, _ := jaeger.InitJaeger("kitex-client-interaction")

		cli, err = interactionservice.NewClient(conf.InteractionService.Name, client.WithHostPorts(conf.InteractionService.Addr), client.WithSuite(tracerSuite))
		if err != nil {
			panic(err)
		}
	})
	return cli, err
}
