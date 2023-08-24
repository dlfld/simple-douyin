// Package handler /*
package handler

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

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
// @Param token body video.PublishActionRequest true "body"
// @Accept json
// @Produce json
// @Router /douyin/publish/action/ [POST]
func VideoSubmit(c *gin.Context) {
	// 创建客户端链接
	cli, err := videoRpc.NewRpcVideoClient()
	if err != nil {
		panic(err)
	}
	// Token: c.Query("token"),
	// 创建发生消息的请求实例 接收视频投稿信息

	req := video.NewPublishActionRequest()
	req.Title = c.PostForm("title")
	req.UserId = int64(c.GetUint("userID"))
	log.Printf("userId in gin = %v\n", req.UserId)
	req.UserId = int64(c.GetUint("userID"))
	req.Data = []byte(c.PostForm("data"))
	file, err := c.FormFile("data")
	fileContent, _ := file.Open()
	byteContainer, err := io.ReadAll(fileContent)
	req.Data = byteContainer

	// 前端请求数据绑定到req中
	// 发起RPC调用
	resp, err := cli.PublishAction(context.Background(), req)

	if err != nil {
		panic(err)
	}
	//返回给前端
	c.JSON(http.StatusOK, resp)
}

// VideoFeed @Summary 视频操作
// @Schemes
// @Description 获取最近新发的30条视频
// @Tags 视频接口
// @Param latest_time query int64 true "可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间"
// @Param token query string true "用户鉴权token"
// @Accept json
// @Produce json
// @Router /douyin/feed [GET]
func VideoFeed(c *gin.Context) {
	// 创建客户端链接
	cli, err := videoRpc.NewRpcVideoClient()
	if err != nil {
		panic(err)
	}
	//token := c.Query("token")
	latestTime := c.Query("latest_time")
	var timestamp int64 = 0
	if latestTime != "" {
		timestamp, _ = strconv.ParseInt(latestTime, 10, 64)
	} else {
		timestamp = time.Now().UnixMilli()
	}
	//封装feed
	feedRequest := &video.FeedRequest{}
	feedRequest.LatestTime = &timestamp
	feedRequest.UserId = int64(c.GetUint("userID"))
	resp, err := cli.Feed(context.Background(), feedRequest)
	if err != nil {
		panic(err)
	}
	//返回给前端
	c.JSON(http.StatusOK, resp)
}
