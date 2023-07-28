package handler

import (
	"context"
	"github.com/douyin/kitex_gen/message"
	"github.com/douyin/rpcClient/messageRpc"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MessageAction(c *gin.Context) {
	// 1. 创建客户端连接
	cli, err := messageRpc.NewRpcMessageClient()
	if err != nil {
		panic(err)
	}
	// 2. 创建发生消息的请求实例
	req := message.NewMessageActionRequest()
	// 3. 前端请求数据绑定到req中
	_ = c.ShouldBind(req)
	// 4. 发起RPC调用
	resp, err := cli.SendMessage(context.Background(), req)
	if err != nil {
		panic(err)
	}
	// 5. gin返回给前端
	c.JSON(http.StatusOK, resp)
}
