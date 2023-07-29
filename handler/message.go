package handler

import (
	"context"
	"net/http"

	"github.com/douyin/kitex_gen/message"
	"github.com/douyin/rpcClient/messageRpc"
	"github.com/gin-gonic/gin"
)

// MessageAction @Summary 消息操作
// @Schemes
// @Description 登录用户对消息的相关操作，目前只支持消息发送
// @Tags example
// @Accept json
// @Produce json
// @Param token body message.MessageActionRequest true "Message Action Params"
// @Router /douyin/message/action/ [POST]
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
