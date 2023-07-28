package relationRpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/douyin/common/conf"
	"github.com/douyin/kitex_gen/relation/relationservice"
)

func NewRpcRelationClient() (relationservice.Client, error) {
	cli, err := relationservice.NewClient(conf.RelationService.Name, client.WithHostPorts(conf.RelationService.Addr))
	if err != nil {
		return nil, err
	}
	return cli, nil
}
