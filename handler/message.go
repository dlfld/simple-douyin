package handler

import (
	"fmt"
	"github.com/douyin/kitex_gen/message"
	"github.com/douyin/rpcClient/messageRpc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	temUid, _ := strconv.Atoi(c.Query("userId"))
	uid := int64(temUid)
	req.FromUserId = uid
	// 4. 发起RPC调用
	resp, err := cli.SendMessage(c, req)
	// 5. 处理错误返回响应
	if resp == nil {
		resp = message.NewMessageActionResponse()
	}
	HandlerErr(resp, err)
	c.JSON(http.StatusOK, resp)
}

func MessageChat(c *gin.Context) {
	// 1. 创建客户端连接
	cli, err := messageRpc.NewRpcMessageClient()
	if err != nil {
		panic(err)
	}
	// 2. 创建请求实例
	req := message.NewMessageChatRequest()
	// 3. 前端请求数据绑定到req中
	fmt.Println(c.Query("userId"))
	uid, _ := strconv.Atoi(c.Query("userId"))
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
