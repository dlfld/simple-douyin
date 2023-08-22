package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/douyin/kitex_gen/interaction"
	"github.com/douyin/kitex_gen/interaction/interactionservice"
	"github.com/douyin/rpcClient/interactionRpc"
	"github.com/gin-gonic/gin"
)

var cli interactionservice.Client
var once sync.Once

// InitInteractionCli 创建一个rpc client 连接
func InitInteractionCli() (err error) {
	once.Do(func() {
		cli, err = interactionRpc.NewRpcInteractionClient()
	})
	if err != nil {
		log.Printf("初始化interaction rpcclient 失败： %+v\n", err)
	}
	return
}

// @Summary xxx
// @Schemes
// @Description xxx
// @Tags 互动接口1
// @Accept json
// @Produce json
// @Param request_body body interaction.FavoriteActionRequest true "request body"
// @Router /douyin/favorite/action/ [POST]
func InteractionFavoriteAction(c *gin.Context) {

	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actionType, err := strconv.Atoi(c.Query("action_type")) // 1-点赞，2-取消点赞
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	req := &interaction.FavoriteActionRequest{
		VideoId:    videoId,
		ActionType: int32(actionType),
		UserId:     userId,
	}

	resp, err := cli.FavoriteAction(context.Background(), req)
	if err != nil || resp.StatusCode != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": resp.StatusMsg,
			"err": err,
		})
	}
	c.JSON(http.StatusOK, resp)

}

// @Summary xxx
// @Schemes
// @Description xxx
// @Tags 互动接口
// @Accept json
// @Produce json
// @Param request_body body interaction.FavoriteListRequest true "request body"
// @Router /douyin/favorite/list/ [GET]
func InteractionFavoriteList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	req := &interaction.FavoriteListRequest{
		UserId: userId,
	}

	resp, err := cli.FavoriteList(context.Background(), req)
	if err != nil || resp.StatusCode != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": resp.StatusMsg,
			"err": err,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary xxx
// @Schemes
// @Description xxx
// @Tags 互动接口
// @Accept json
// @Produce json
// @Param request_body body interaction.CommentActionRequest true "request body"
// @Router /douyin/comment/action/ [POST]
func InteractionCommentAction(c *gin.Context) {

	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actionType, err := strconv.Atoi(c.Query("action_type")) // 1-点赞，2-取消点赞
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 32)
	commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
	commentText := c.Query("comment_text")
	req := &interaction.CommentActionRequest{
		VideoId:     videoId,
		UserId:      &userId,
		ActionType:  int32(actionType),
		CommentText: &commentText,
		CommentId:   &commentId,
	}

	resp, err := cli.CommentAction(context.Background(), req)
	if err != nil || resp.StatusCode != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": resp.StatusMsg,
			"err": err,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary xxx
// @Schemes
// @Description xxx
// @Tags 互动接口
// @Accept json
// @Produce json
// @Param request_body body interaction.CommentListRequest true "request body"
// @Router /douyin/comment/list/ [GET]
func InteractionCommentList(c *gin.Context) {

	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	req := &interaction.CommentListRequest{
		UserId:  &userId,
		VideoId: videoId,
	}

	resp, err := cli.CommentList(context.Background(), req)
	if err != nil || resp.StatusCode != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": resp.StatusMsg,
			"err": err,
		})
	}

	c.JSON(http.StatusOK, resp)
}
