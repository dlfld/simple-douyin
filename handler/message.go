package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/douyin/kitex_gen/message"
	"github.com/douyin/rpcClient/messageRpc"
	"github.com/gin-gonic/gin"
)

// MessageAction @Summary 消息操作
// @Schemes
// @Description 登录用户对消息的相关操作，目前只支持消息发送
// @Tags 消息
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
	fmt.Println(req)
	temUid, _ := strconv.Atoi(c.Query("userID"))
	// uid := int64(temUid)
	// user := model.BaseReq{
	// 	UserId: &uid,
	// }
	// req.User = &user
	req.FromUserId = int64(temUid)
	// 4. 发起RPC调用
	//ctx := context.WithValue(context.Background(), "userID", c.Query(""))
	resp, err := cli.SendMessage(c, req)
	if err != nil {
		c.JSON(http.StatusOK, resp)
	}
	// 5. gin返回给前端
	c.JSON(http.StatusOK, resp)
}

// MessageAction @Summary 消息操作
// @Schemes
// @Description 登录用户对消息的相关操作，目前只支持消息发送
// @Tags 消息
// @Accept json
// @Produce json
// @Param token body message.MessageChatRequest true "Message Action Params"
// @Router /douyin/message/chat/ [POST]
func MessageChat(c *gin.Context) {
	// 1. 创建客户端连接
	cli, err := messageRpc.NewRpcMessageClient()
	if err != nil {
		panic(err)
	}
	// 2. 创建请求实例
	req := message.NewMessageChatRequest()
	// 3. 前端请求数据绑定到req中
	_ = c.ShouldBind(req)
	fmt.Println(req)
	temUid, _ := strconv.Atoi(c.Query("userID"))
	// uid := int64(temUid)
	// user := model.BaseReq{
	// 	UserId: &uid,
	// }
	// req.User = &user
	req.FromUserId = int64(temUid)
	// 4. 发起RPC调用
	//ctx := context.WithValue(context.Background(), "userID", c.Query(""))
	resp, err := cli.MessageList(c, req)
	if err != nil {
		c.JSON(http.StatusOK, resp)
	}
	// 5. gin返回给前端
	c.JSON(http.StatusOK, resp)
}
