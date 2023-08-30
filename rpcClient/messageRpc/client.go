package messageRpc

import (
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/jaeger"
	"github.com/douyin/kitex_gen/message/messageservice"
)

var once sync.Once
var cli messageservice.Client
var err error

func NewRpcMessageClient() (messageservice.Client, error) {
	once.Do(func() {
		tracerSuite, _ := jaeger.InitJaeger("kitex-client-message")
		cli, err = messageservice.NewClient(conf.MessageService.Name, client.WithHostPorts(conf.MessageService.Addr), client.WithSuite(tracerSuite))
	})
	return cli, err
}
