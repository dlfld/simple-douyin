// Package handler /*
package handler

import (
	"context"
	"net/http"

	"github.com/douyin/kitex_gen/video"
	"github.com/douyin/rpcClient/videoRpc"
	"github.com/gin-gonic/gin"
)

// PublishList @Summary 视频操作
// @Schemes
// @Description 获取用户发表的视频列表
// @Tags 视频接口
// @Param token query video.PublishListRequest true "body"
// @Accept json
// @Produce json
// @Router /douyin/publish/list/ [GET]
func PublishList(c *gin.Context) {
	// 创建客户端链接
	cli, err := videoRpc.NewRpcVideoClient()
	if err != nil {
		panic(err)
	}
	// 创建发生消息的请求实例
	req := video.NewPublishListRequest()
	// 前端请求数据绑定到req中
	_ = c.ShouldBind(req)
	// 发起RPC调用
	resp, err := cli.PublishList(context.Background(), req)
	if err != nil {
		panic(err)
	}
	//返回给前端
	c.JSON(http.StatusOK, resp)
}
