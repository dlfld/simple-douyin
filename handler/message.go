package handler

import (
	"fmt"
	"github.com/douyin/common/constant"
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
		c.JSON(http.StatusOK, constant.NewErrResp(constant.ErrSystemBusy))
		return
	}
	// 2. 创建发生消息的请求实例
	req := message.NewMessageActionRequest()
	// 3. 前端请求数据绑定到req中
	err = c.ShouldBind(req)
	if req.Content == "" {
		c.JSON(http.StatusOK, constant.NewErrResp(constant.ErrEmptyMessage))
		return
	}
	if len(req.Content) > 1000 {
		c.JSON(http.StatusOK, constant.NewErrResp(constant.ErrMsgTooLong))
		return
	}
	req.FromUserId = int64(c.GetUint("userID"))
	fmt.Printf("req = %+v\n", req)
	// 4. 发起RPC调用
	resp, err := cli.SendMessage(c, req)
	// 5. 处理错误返回响应
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
		c.JSON(http.StatusOK, constant.NewErrResp(constant.ErrSystemBusy))
		return
	}
	// 2. 创建请求实例
	req := message.NewMessageChatRequest()
	// 3. 前端请求数据绑定到req中
	uid := c.GetUint("userID")
	toUserID, _ := strconv.Atoi(c.Query("to_user_id"))
	timeStr := c.Query("pre_msg_time")
	preMsgTime, _ := strconv.Atoi(timeStr)
	req.FromUserId = int64(uid)
	req.ToUserId = int64(toUserID)
	req.PreMsgTime = int64(preMsgTime)
	// 4. 发起RPC调用
	resp, err := cli.MessageList(c, req)
	// 5. 处理错误返回响应
	c.JSON(http.StatusOK, resp)
}
