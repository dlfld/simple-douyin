// Package handler /*
package handler

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/douyin/common/constant"
	"github.com/douyin/kitex_gen/video"
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
	//cli, err := videoRpc.NewRpcVideoClient()
	//if err != nil {
	//	panic(err)
	//}
	// 创建发生消息的请求实例
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, constant.NewErrResp(constant.ErrBadRequest))
		return
	}
	req := &video.PublishListRequest{
		UserId: int64(userId),
		Token:  c.Query("token"),
	}

	// 发起RPC调用
	resp, err := rpcCli.videoCli.PublishList(context.Background(), req)
	if err != nil {
		constant.HandlerErr(constant.ErrPublishList, resp)
	}
	//返回给前端
	c.JSON(http.StatusOK, resp)
}

// VideoSubmit @Summary 视频操作
// @Schemes
// @Description 视频投稿
// @Tags 视频接口
// @Param data formData file true "file"
// @Param title formData string true "title"
// @Param user_id query string true "user_id"
// @Param token query string true "token"
// @Accept json
// @Produce json
// @Router /douyin/publish/action/ [POST]
func VideoSubmit(c *gin.Context) {

	req := video.NewPublishActionRequest()
	req.Title = c.PostForm("title")
	req.UserId = int64(c.GetUint("userID"))
	log.Printf("userId in gin = %v\n", req.UserId)
	req.UserId = int64(c.GetUint("userID"))
	req.Data = []byte(c.PostForm("data"))
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, constant.NewErrResp(constant.ErrBadRequest))
		return
	}
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, constant.NewErrResp(constant.ErrFileCantOpen))
		return
	}
	byteContainer, err := io.ReadAll(fileContent)
	req.Data = byteContainer
	if err != nil {
		c.JSON(http.StatusOK, constant.NewErrResp(constant.ErrBadRequest))
		return
	}
	// 前端请求数据绑定到req中
	// 发起RPC调用
	resp, err := rpcCli.videoCli.PublishAction(context.Background(), req)

	if err != nil {
		constant.HandlerErr(constant.ErrVideoPublish, resp)
	}
	//返回给前端
	c.JSON(http.StatusOK, resp)
}

// VideoFeed @Summary 视频操作
// @Schemes
// @Description 获取最近新发的30条视频
// @Tags 视频接口
// @Param latest_time query int64 false "可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间"
// @Param token query string true "用户鉴权token"
// @Accept json
// @Produce json
// @Router /douyin/feed [GET]
func VideoFeed(c *gin.Context) {
	// 创建客户端链接

	latestTime := c.Query("latest_time")
	var timestamp int64 = 0
	var err error
	if latestTime != "" {
		timestamp, err = strconv.ParseInt(latestTime, 10, 64)
	} else {
		timestamp = time.Now().UnixMilli()
	}
	if err != nil {
		c.JSON(http.StatusOK, constant.NewErrResp(constant.ErrBadRequest))
		return
	}
	//封装feed
	feedRequest := &video.FeedRequest{}
	feedRequest.LatestTime = &timestamp
	feedRequest.UserId = int64(c.GetUint("userID"))
	resp, err := rpcCli.videoCli.Feed(context.Background(), feedRequest)
	if err != nil {
		constant.HandlerErr(constant.ErrFeedErr, resp)
		return
	}
	//返回给前端
	c.JSON(http.StatusOK, resp)
}
