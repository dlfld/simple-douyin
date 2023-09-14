package messageRpc

import (
	"github.com/douyin/common/etcd"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/kitex_gen/message/messageservice"
)

var once sync.Once
var cli messageservice.Client
var err error

func NewRpcMessageClient() (messageservice.Client, error) {
	once.Do(func() {
		addr := etcd.DiscoverService(conf.MessageService.Name)
		cli, err = messageservice.NewClient(conf.MessageService.Name, client.WithHostPorts(addr...), client.WithSuite(tracing.NewClientSuite()))
	})
	return cli, err
}
