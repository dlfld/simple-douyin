package messageRpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/kitex_gen/message/messageservice"
)

func NewRpcMessageClient() (messageservice.Client, error) {
	cli, err := messageservice.NewClient(conf.MessageService.Name, client.WithHostPorts(conf.MessageService.Addr))
	if err != nil {
		return nil, err
	}
	return cli, nil
}
