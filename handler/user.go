package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/douyin/kitex_gen/user"
	"github.com/gin-gonic/gin"
)

// Register @Summary 用户注册
// @Schemes
// @Description 初始注册。
// @Tags 基础接口
// @Accept json
// @Produce json
// @Param username query string true "注册昵称"
// @Param password query string true "注册密码"
// @Router /douyin/user/register/  [POST]
func Register(c *gin.Context) {
	// 1. 创建客户端连接
	//cli, err := userRpc.NewRpcUserClient()
	//if err != nil {
	//	panic(err)
	//}
	// 2. 创建发生消息的请求实例
	// req := relation.NewFollowerListRequest()
	username := c.Query("username")
	password := c.Query("password")
	req := &user.UserRegisterRequest{Username: username, Password: password}
	// 3. 前端请求数据绑定到req中
	// _ = c.ShouldBind(req)
	// 4. 发起RPC调用
	resp, _ := rpcCli.userCli.UserRegister(context.Background(), req)
	//if err2 != nil {
	//	panic(err2)
	//}
	// 5. gin返回给前端
	c.JSON(http.StatusOK, resp)
}

// Login @Summary 用户登录
// @Schemes
// @Description 登录。
// @Tags 基础接口
// @Accept json
// @Produce json
// @Param username query string true "注册昵称"
// @Param password query string true "注册密码"
// @Router /douyin/user/login/  [POST]
func Login(c *gin.Context) {
	// 1. 创建客户端连接
	//cli, err := userRpc.NewRpcUserClient()
	//if err != nil {
	//	panic(err)
	//}
	// 2. 创建发生消息的请求实例
	// req := relation.NewFollowerListRequest()
	username := c.Query("username")
	password := c.Query("password")
	req := &user.UserLoginRequest{Username: username, Password: password}
	// 3. 前端请求数据绑定到req中
	// _ = c.ShouldBind(req)
	// 4. 发起RPC调用
	resp, err2 := rpcCli.userCli.UserLogin(context.Background(), req)
	if err2 != nil {
		c.JSON(http.StatusBadGateway, resp)
		//panic(err2)
	}
	// 5. gin返回给前端
	c.JSON(http.StatusOK, resp)
}

// UserInfo @Summary 用户信息
// @Schemes
// @Description 获取用户基础信息。
// @Tags 基础接口
// @Accept json
// @Produce json
// @Param user_id query int true "用户id"
// @Param token query string true "用户鉴权token"
// @Router /douyin/user/ [GET]
func UserInfo(c *gin.Context) {
	// 1. 创建客户端连接
	//cli, err := userRpc.NewRpcUserClient()
	//if err != nil {
	//	panic(err)
	//}
	// 2. 创建发生消息的请求实例
	// req := relation.NewFollowerListRequest()
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	//if err2 != nil {
	//	panic(err)
	//}
	token := c.Query("token")
	req := &user.UserRequest{UserId: userId, Token: token}
	// 3. 前端请求数据绑定到req中
	// _ = c.ShouldBind(req)
	// 4. 发起RPC调用
	resp, _ := rpcCli.userCli.UserMsg(context.Background(), req)
	if resp.User != nil {
		c.Set("userId", resp.User.Id)
	}
	//if err3 != nil {
	//	panic(err)
	//}
	// 5. gin返回给前端
	c.JSON(http.StatusOK, resp)
}
