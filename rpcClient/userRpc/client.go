package userRpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/kitex_gen/user/userservice"
)

func NewRpcUserClient() (userservice.Client, error) {
	cli, err := userservice.NewClient(conf.UserService.Name, client.WithHostPorts(conf.UserService.Addr))
	if err != nil {
		return nil, err
	}
	return cli, nil
}
