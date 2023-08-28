package userRpc

import (
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/jaeger"
	"github.com/douyin/kitex_gen/user/userservice"
)

var cli userservice.Client
var once sync.Once
var err error

func NewRpcUserClient() (userservice.Client, error) {
	once.Do(func() {
		tracerSuite, _ := jaeger.InitJaeger("kitex-client-user")
		cli, err = userservice.NewClient(conf.UserService.Name, client.WithHostPorts(conf.UserService.Addr), client.WithSuite(tracerSuite))
		if err != nil {
			panic(err)
		}
	})
	return cli, nil
}
