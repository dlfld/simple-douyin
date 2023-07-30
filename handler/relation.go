package handler

import (
	"context"
	"net/http"

	"github.com/douyin/kitex_gen/relation"
	"github.com/douyin/rpcClient/relationRpc"
	"github.com/gin-gonic/gin"
)

// @Summary 关系操作
// @Schemes
// @Description 登录用户对其他用户进行关注或取消关注。
// @Accept json
// @Produce json
// @Param token body relation.FollowActionRequest true "request body"
// @Router /douyin/relation/action/ [POST]
func RelationAction(c *gin.Context) {
	// 1. 创建客户端连接
	cli, err := relationRpc.NewRpcRelationClient()
	if err != nil {
		panic(err)
	}
	// 2. 创建发生消息的请求实例
	req := relation.NewFollowActionRequest()
	// 3. 前端请求数据绑定到req中
	_ = c.ShouldBind(req)
	// 4. 发起RPC调用
	resp, err := cli.FollowAction(context.Background(), req)
	if err != nil {
		panic(err)
	}
	// 5. gin返回给前端
	c.JSON(http.StatusOK, resp)
}
