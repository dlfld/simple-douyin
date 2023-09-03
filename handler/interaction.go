package handler

import (
	"context"
	"strconv"

	"github.com/douyin/kitex_gen/interaction"

	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 用户点赞操作
// @Schemes
// @Description 登录用户对视频的点赞和取消点赞操作。
// @Tags 互动接口
// @Accept json
// @Produce json
// @Param request_body query interaction.FavoriteActionRequest true "request body"
// @Router /douyin/favorite/action/ [POST]
func InteractionFavoriteAction(c *gin.Context) {

	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actionType, err := strconv.Atoi(c.Query("action_type")) // 1-点赞，2-取消点赞
	userId := int64(c.GetUint("userID"))
	req := &interaction.FavoriteActionRequest{
		VideoId:    videoId,
		ActionType: int32(actionType),
		UserId:     int64(userId),
	}

	resp, err := rpcCli.interactionCli.FavoriteAction(context.Background(), req)
	if err != nil || resp.StatusCode != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": resp.StatusMsg,
			"err": err,
		})
		return
	}
	c.JSON(http.StatusOK, resp)

}

// @Summary 获取用户点赞列表
// @Schemes
// @Description 登录用户的所有点赞视频。
// @Tags 互动接口
// @Accept json
// @Produce json
// @Param request_body query interaction.FavoriteListRequest true "request body"
// @Router /douyin/favorite/list/ [GET]
func InteractionFavoriteList(c *gin.Context) {
	// userId := int64(c.GetUint("userID"))
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		return
	}
	req := &interaction.FavoriteListRequest{
		UserId: int64(userId),
	}

	resp, err := rpcCli.interactionCli.FavoriteList(context.Background(), req)
	// log.Printf("resp = %+v\n", resp)
	if resp != nil && err != nil || resp.StatusCode != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": resp.StatusMsg,
			"err": err,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary 用户评论操作
// @Schemes
// @Description 登录用户对视频进行评论。
// @Tags 互动接口
// @Accept json
// @Produce json
// @Param request_body query interaction.CommentActionRequest true "request body"
// @Router /douyin/comment/action/ [POST]
func InteractionCommentAction(c *gin.Context) {

	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actionType, err := strconv.Atoi(c.Query("action_type")) // 1-点赞，2-取消点赞
	userId := int64(c.GetUint("userID"))
	commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
	commentText := c.Query("comment_text")
	req := &interaction.CommentActionRequest{
		VideoId:     videoId,
		UserId:      &userId,
		ActionType:  int32(actionType),
		CommentText: &commentText,
		CommentId:   &commentId,
	}
	resp, err := rpcCli.interactionCli.CommentAction(context.Background(), req)
	if err != nil || resp.StatusCode != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": resp.StatusMsg,
			"err": err,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary 获取视频评论
// @Schemes
// @Description 查看视频的所有评论，按发布时间倒序。
// @Tags 互动接口
// @Accept json
// @Produce json
// @Param request_body query interaction.CommentListRequest true "request body"
// @Router /douyin/comment/list/ [GET]
func InteractionCommentList(c *gin.Context) {

	userId := int64(c.GetUint("userID"))
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	req := &interaction.CommentListRequest{
		UserId:  &userId,
		VideoId: videoId,
	}

	resp, err := rpcCli.interactionCli.CommentList(context.Background(), req)
	if err != nil || resp.StatusCode != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": resp.StatusMsg,
			"err": err,
		})
	}

	c.JSON(http.StatusOK, resp)
}
