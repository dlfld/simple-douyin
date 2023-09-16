package videoRpc

import (
	"github.com/douyin/common/etcd"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/kitex_gen/video/videoservice"
)

var cli videoservice.Client
var once sync.Once
var err error

func NewRpcVideoClient() (videoservice.Client, error) {
	once.Do(func() {
		addr := etcd.DiscoverService(conf.VideoService.Name)
		cli, err = videoservice.NewClient(conf.VideoService.Name, client.WithHostPorts(addr...), client.WithSuite(tracing.NewClientSuite()))
		if err != nil {
			panic(err)
		}
	})
	return cli, nil
}
