package videoRpc

import (
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/jaeger"
	"github.com/douyin/kitex_gen/video/videoservice"
)

var cli videoservice.Client
var once sync.Once
var err error

func NewRpcVideoClient() (videoservice.Client, error) {
	once.Do(func() {
		tracerSuite, _ := jaeger.InitJaeger("kitex-client-video")
		cli, err = videoservice.NewClient(conf.VideoService.Name, client.WithHostPorts(conf.VideoService.Addr), client.WithSuite(tracerSuite))
		if err != nil {
			panic(err)
		}
	})
	return cli, nil
}
