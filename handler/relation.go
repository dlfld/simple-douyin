package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/douyin/kitex_gen/relation"
	"github.com/douyin/rpcClient/relationRpc"
	"github.com/gin-gonic/gin"
)

// var relationCliPool = sync.Pool{
// 	New: func() any {
// 		cli, _ := relationRpc.NewRpcRelationClient()
// 		return cli
// 	},
// }

// @Summary 关系操作
// @Schemes
// @Description 登录用户对其他用户进行关注或取消关注。
// @Tags 社交接口
// @Accept json
// @Produce json
// @Param request_body body relation.FollowActionRequest true "request body"
// @Router /douyin/relation/action/ [POST]
func RelationAction(c *gin.Context) {
	// 1. 创建客户端连接
	// cli := relationCliPool.Get().(relationservice.Client)
	cli, err := relationRpc.NewRpcRelationClient()

	if err != nil {
		panic(err)
	}
	// 2. 创建发生消息的请求实例
	req := relation.NewFollowActionRequest()
	// 3. 前端请求数据绑定到req中
	err = c.ShouldBindJSON(req)
	if err != nil {
		panic(err)
	}
	// 4. 发起RPC调用
	resp, err := cli.FollowAction(context.Background(), req)
	if err != nil {
		panic(err)
	}
	// defer relationCliPool.Put(cli)
	// 5. gin返回给前端
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
	// 1. 创建客户端连接
	cli, err := relationRpc.NewRpcRelationClient()
	if err != nil {
		panic(err)
	}
	// cli := relationCliPool.Get().(relationservice.Client)
	// 2. 创建发生消息的请求实例
	// req := relation.NewFollowingListRequest()
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		panic(err)
	}
	req := &relation.FollowingListRequest{
		UserId: int64(userId),
		Token:  c.Query("token"),
	}
	// 3. 前端请求数据绑定到req中
	// _ = c.ShouldBindQuery(req)

	// 4. 发起RPC调用
	resp, err := cli.FollowList(context.Background(), req)
	if err != nil {
		panic(err)
	}
	// relationCliPool.Put(cli)
	// 5. gin返回给前端
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
	// 1. 创建客户端连接
	cli, err := relationRpc.NewRpcRelationClient()
	if err != nil {
		panic(err)
	}
	// 2. 创建发生消息的请求实例
	// req := relation.NewFollowerListRequest()
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		panic(err)
	}
	req := &relation.FollowerListRequest{
		UserId: int64(userId),
		Token:  c.Query("token"),
	}
	// 3. 前端请求数据绑定到req中
	// _ = c.ShouldBind(req)
	// 4. 发起RPC调用
	resp, err := cli.FollowerList(context.Background(), req)
	if err != nil {
		panic(err)
	}
	// 5. gin返回给前端
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
	// 1. 创建客户端连接
	cli, err := relationRpc.NewRpcRelationClient()
	if err != nil {
		panic(err)
	}
	// 2. 创建发生消息的请求实例
	// req := relation.NewRelationFriendListRequest()
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		panic(err)
	}
	req := &relation.RelationFriendListRequest{
		UserId: int64(userId),
		Token:  c.Query("token"),
	}
	// 3. 前端请求数据绑定到req中
	// _ = c.ShouldBind(req)
	// 4. 发起RPC调用
	resp, err := cli.FriendList(context.Background(), req)
	if err != nil {
		panic(err)
	}
	// 5. gin返回给前端
	c.JSON(http.StatusOK, resp)
}
