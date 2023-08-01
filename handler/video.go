// Package handler /*
package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/douyin/kitex_gen/video"
	"github.com/douyin/rpcClient/videoRpc"
	"github.com/gin-gonic/gin"
)

// PublishList @Summary 视频操作
// @Schemes
// @Description 获取用户发表的视频列表
// @Tags 视频接口
// @Param token body video.PublishListRequest true "body"
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
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		panic(err)
	}
	req := &video.PublishListRequest{
		UserId: int64(userId),
		Token:  c.Query("token"),
	}
	// 发起RPC调用
	resp, err := cli.PublishList(context.Background(), req)
	if err != nil {
		panic(err)
	}
	//返回给前端
	c.JSON(http.StatusOK, resp)
}

// VideoSubmit @Summary 视频操作
// @Schemes
// @Description 视频投稿
// @Tags 视频接口
// @Param token query video.PublishActionRequest true "body"
// @Accept json
// @Produce json
// @Router /douyin/publish/action/ [POST]
func VideoSubmit(c *gin.Context) {
	// 创建客户端链接
	cli, err := videoRpc.NewRpcVideoClient()
	if err != nil {
		panic(err)
	}
	// 创建发生消息的请求实例 接收视频投稿信息
	req := video.NewPublishActionRequest()

	if err != nil {
		return
	}
	// 前端请求数据绑定到req中
	_ = c.ShouldBind(req)

	// 发起RPC调用
	resp, err := cli.PublishAction(context.Background(), req)

	if err != nil {
		panic(err)
	}
	//返回给前端
	c.JSON(http.StatusOK, resp)
}
