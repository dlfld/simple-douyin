package handler

import (
	"fmt"
	"github.com/douyin/kitex_gen/message"
	"github.com/douyin/kitex_gen/model"
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
	fmt.Println(req)
	temUid, _ := strconv.Atoi(c.Query("userID"))
	uid := int64(temUid)
	user := model.BaseReq{
		UserId: &uid,
	}
	req.User = &user
	// 4. 发起RPC调用
	//ctx := context.WithValue(context.Background(), "userID", c.Query(""))
	resp, err := cli.SendMessage(c, req)
	if err != nil {
		c.JSON(http.StatusOK, resp)
	}
	// 5. gin返回给前端
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
	_ = c.ShouldBind(req)
	fmt.Println(req)
	temUid, _ := strconv.Atoi(c.Query("userID"))
	uid := int64(temUid)
	user := model.BaseReq{
		UserId: &uid,
	}
	req.User = &user
	// 4. 发起RPC调用
	//ctx := context.WithValue(context.Background(), "userID", c.Query(""))
	resp, err := cli.MessageList(c, req)
	if err != nil {
		c.JSON(http.StatusOK, resp)
	}
	// 5. gin返回给前端
	c.JSON(http.StatusOK, resp)
}
