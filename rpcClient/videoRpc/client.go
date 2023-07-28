package videoRpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/kitex_gen/video/videoservice"
)

func NewRpcVideoClient() (videoservice.Client, error) {
	cli, err := videoservice.NewClient(conf.VideoService.Name, client.WithHostPorts(conf.VideoService.Addr))
	if err != nil {
		return nil, err
	}
	return cli, nil
}
