package messageRpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/kitex_gen/message/messageservice"
	"sync"
)

var once sync.Once
var cli messageservice.Client
var err error

func NewRpcMessageClient() (messageservice.Client, error) {
	once.Do(func() {
		cli, err = messageservice.NewClient(conf.MessageService.Name, client.WithHostPorts(conf.MessageService.Addr))
	})
	return cli, err
}
