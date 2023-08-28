package handler

import (
	"log"
	"sync"

	"github.com/douyin/kitex_gen/interaction/interactionservice"
	"github.com/douyin/kitex_gen/message/messageservice"
	"github.com/douyin/kitex_gen/relation/relationservice"
	"github.com/douyin/kitex_gen/user/userservice"
	"github.com/douyin/kitex_gen/video/videoservice"
	"github.com/douyin/rpcClient/interactionRpc"
	"github.com/douyin/rpcClient/messageRpc"
	"github.com/douyin/rpcClient/relationRpc"
	"github.com/douyin/rpcClient/userRpc"
	"github.com/douyin/rpcClient/videoRpc"
)

var rpcCli *rpcClients
var once sync.Once

type rpcClients struct {
	interactionCli interactionservice.Client
	messageCli     messageservice.Client
	relationCli    relationservice.Client
	userCli        userservice.Client
	videoCli       videoservice.Client
}

func InitRpcCli() {
	//todo err做成list发送到kafka
	var err error
	once.Do(func() {
		rpcCli.interactionCli, err = interactionRpc.NewRpcInteractionClient()
		if err != nil {
			log.Printf("初始化interaction rpcclient 失败： %+v\n", err)
		}
		rpcCli.messageCli, err = messageRpc.NewRpcMessageClient()
		if err != nil {
			log.Printf("初始化message rpcclient 失败： %+v\n", err)
		}
		rpcCli.relationCli, err = relationRpc.NewRpcRelationClient()
		if err != nil {
			log.Printf("初始化relation rpcclient 失败： %+v\n", err)
		}
		rpcCli.userCli, err = userRpc.NewRpcUserClient()
		if err != nil {
			log.Printf("初始化user rpcclient 失败： %+v\n", err)
		}
		rpcCli.videoCli, err = videoRpc.NewRpcVideoClient()
		if err != nil {
			log.Printf("初始化video rpcclient 失败： %+v\n", err)
		}
	})
}
