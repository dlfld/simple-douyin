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
	_ = c.ShouldBind(&req)
	atoi, err := strconv.Atoi(c.Query("action_type"))
	req.ActionType = int32(atoi)
	toUserId, err := strconv.Atoi(c.Query("to_user_id"))
	req.ToUserId = int64(toUserId)
	req.Content = c.Query("content")
	//temUid, _ := strconv.Atoi(c.Query("userId"))
	temUid := c.GetUint("userID")
	uid := int64(temUid)
	req.FromUserId = uid
	fmt.Printf("req = %+v\n", req)
	// 4. 发起RPC调用
	resp, err := cli.SendMessage(c, req)
	// 5. 处理错误返回响应
	if resp == nil {
		resp = message.NewMessageActionResponse()
	}
	HandlerErr(resp, err)
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
	fmt.Println("userid")
	fmt.Println(c.GetUint("userID"))

	uid := c.GetUint("userID")
	toUserID, _ := strconv.Atoi(c.Query("to_user_id"))
	preMsgTime, _ := strconv.Atoi(c.Query("pre_msg_time"))
	req.FromUserId = int64(uid)
	req.ToUserId = int64(toUserID)
	req.PreMsgTime = int64(preMsgTime)
	// 4. 发起RPC调用
	resp, err := cli.MessageList(c, req)
	if resp == nil {
		resp = message.NewMessageChatResponse()
	}
	// 5. 处理错误返回响应
	HandlerErr(resp, err)
	c.JSON(http.StatusOK, resp)
}
