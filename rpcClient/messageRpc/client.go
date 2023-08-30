package messageRpc

import (
	"github.com/douyin/common/etcd"
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
		addr := etcd.DiscoverService(conf.MessageService.Name)
		cli, err = messageservice.NewClient(conf.MessageService.Name, client.WithHostPorts(addr...), client.WithSuite(tracerSuite))
	})
	return cli, err
}
