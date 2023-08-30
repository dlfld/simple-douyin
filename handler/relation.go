package handler

import (
	"context"
	"fmt"

	"net/http"
	"strconv"

	"github.com/douyin/common/constant"
	"github.com/douyin/kitex_gen/relation"
	"github.com/douyin/kitex_gen/relation/relationservice"
	"github.com/douyin/rpcClient/relationRpc"
	"github.com/gin-gonic/gin"
)

var relationCli relationservice.Client

func init() {
	var err error
	relationCli, err = relationRpc.NewRpcRelationClient()
	if err != nil {
		panic(err)
	}
}

// @Summary 关系操作
// @Schemes
// @Description 登录用户对其他用户进行关注或取消关注。
// @Tags 社交接口
// @Accept json
// @Produce json
// @Param request_body body relation.FollowActionRequest true "request body"
// @Router /douyin/relation/action/ [POST]
func RelationAction(c *gin.Context) {
	ToUserID, _ := strconv.Atoi(c.Query("to_user_id"))
	ActionType, _ := strconv.Atoi(c.Query("action_type"))
	userID, has := c.Get("userID")
	fmt.Println(userID, has)
	req := relation.FollowActionRequest{
		FromUserId: int64(c.GetUint("userID")),
		ToUserId:   int64(ToUserID),
		ActionType: int32(ActionType),
	}
	fmt.Println("req:", req)

	// 发起RPC调用
	resp, err := rpcCli.relationCli.FollowAction(context.Background(), &req)
	if err != nil {
		constant.HandlerErr(constant.ErrFollowAction, resp)
	}

	// gin返回给前端
	c.JSON(http.StatusOK, resp)
}

// @Summary 用户关注列表
// @Schemes
// @Description 登录用户关注的所有用户列表。
// @Tags 社交接口
// @Accept json
// @Produce json
// @Param user_id query int true "用户id"
// @Param token query string true "用户鉴权token"
// @Router /douyin/relation/follow/list/  [GET]
func RelationFollowList(c *gin.Context) {
	// 创建发生消息的请求实例
	// req := relation.NewFollowingListRequest()
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		resp := relation.FollowingListResponse{}
		constant.HandlerErr(constant.ErrBadRequest, &resp)
		c.JSON(http.StatusOK, resp)
		return
	}
	// 前端请求数据绑定到req中
	req := &relation.FollowingListRequest{
		UserId: int64(userId),
		Token:  c.Query("token"),
	}

	// 发起RPC调用
	resp, err := rpcCli.relationCli.FollowList(context.Background(), req)
	if err != nil {
		constant.HandlerErr(constant.ErrGetFollowList, resp)
	}
	// relationCliPool.Put(cli)
	// gin返回给前端
	c.JSON(http.StatusOK, resp)
}

// @Summary 用户粉丝列表
// @Schemes
// @Description 所有关注登录用户的粉丝列表。
// @Tags 社交接口
// @Accept json
// @Produce json
// @Param user_id query int true "用户id"
// @Param token query string true "用户鉴权token"
// @Router /douyin/relation/follower/list/  [GET]
func RelationFollowerList(c *gin.Context) {

	// 创建发生消息的请求实例
	// req := relation.NewFollowerListRequest()
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		resp := relation.FollowerListResponse{}
		constant.HandlerErr(constant.ErrBadRequest, &resp)
		c.JSON(http.StatusOK, resp)
	}
	// 前端请求数据绑定到req中
	req := &relation.FollowerListRequest{
		UserId: int64(userId),
		Token:  c.Query("token"),
	}
	// 发起RPC调用
	resp, err := rpcCli.relationCli.FollowerList(context.Background(), req)
	if err != nil {
		constant.HandlerErr(constant.ErrGetFollowerList, &resp)
	}
	// gin返回给前端
	c.JSON(http.StatusOK, resp)
}

// @Summary 用户好友列表
// @Schemes
// @Description 所有与登录用户互相关注的用户列表
// @Tags 社交接口
// @Accept json
// @Produce json
// @Param user_id query int true "用户id"
// @Param token query string true "用户鉴权token"
// @Router /douyin/relation/friend/list/  [GET]
func RelationFriendList(c *gin.Context) {
	// 创建发生消息的请求实例
	// req := relation.NewRelationFriendListRequest()
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		resp := relation.FollowingListResponse{}
		constant.HandlerErr(constant.ErrBadRequest, &resp)
		c.JSON(http.StatusOK, resp)
	}
	// 前端请求数据绑定到req中
	req := &relation.RelationFriendListRequest{
		UserId: int64(userId),
	}
	// 发起RPC调用
	resp, err := rpcCli.relationCli.FriendList(context.Background(), req)
	if err != nil {
		constant.HandlerErr(constant.ErrGetFriendList, &resp)
	}
	// gin返回给前端
	c.JSON(http.StatusOK, resp)
}
