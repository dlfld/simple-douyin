package userRpc

import (
	"github.com/douyin/common/etcd"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/kitex_gen/user/userservice"
)

var cli userservice.Client
var once sync.Once
var err error

func NewRpcUserClient() (userservice.Client, error) {
	once.Do(func() {
		addr := etcd.DiscoverService(conf.UserService.Name)
		cli, err = userservice.NewClient(conf.UserService.Name, client.WithHostPorts(addr...), client.WithSuite(tracing.NewClientSuite()))
		if err != nil {
			panic(err)
		}
	})
	return cli, nil
}
