package interactionRpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/kitex_gen/interaction/interactionservice"
)

func NewRpcInteractionClient() (interactionservice.Client, error) {
	cli, err := interactionservice.NewClient(conf.InteractionService.Name, client.WithHostPorts(conf.InteractionService.Addr))
	if err != nil {
		return nil, err
	}
	return cli, nil
}
