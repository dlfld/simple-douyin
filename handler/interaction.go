package handler

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/douyin/kitex_gen/interaction"
	"github.com/douyin/kitex_gen/interaction/interactionservice"
	"github.com/douyin/rpcClient/interactionRpc"
	"github.com/gin-gonic/gin"
)

var cli interactionservice.Client
var once sync.Once

// NewMysqlConn 创建一个db数据库
func initInteractionCli() (err error) {
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
// @Tags 互动接口
// @Accept json
// @Produce json
// @Param request_body body relation.FavoriteActionRequest true "request body"
// @Router /douyin/favorite/action/ [POST]

func InteractionFavoriteAction(c *gin.Context) {
	// 1. 创建客户端连接
	err := initInteractionCli()
	if err != nil {
		panic(err)
	}

	// 2. 创建发生消息的请求实例
	req := interaction.NewFavoriteActionRequest()

	// 3. 前端请求数据绑定到req中
	err = c.ShouldBindJSON(req)
	if err != nil {
		panic(err)
	}

	// 4. 发起RPC调用
	resp, err := cli.FavoriteAction(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		panic(err)
	}

	// 5. gin返回给前端
	c.JSON(http.StatusOK, resp)

}

// @Summary xxx
// @Schemes
// @Description xxx
// @Tags 互动接口
// @Accept json
// @Produce json
// @Param request_body body relation.FavoriteActionRequest true "request body"
// @Router /douyin/favorite/list/ [GET]

func InteractionFavoriteList(c *gin.Context) {
	//c.JSON(http.StatusOK, gin.H{"msg": "ok"})

	// 1. 创建客户端连接
	err := initInteractionCli()
	if err != nil {
		panic(err)
	}

	// 2. 创建发生消息的请求实例
	// 3. 前端请求数据绑定到req中
	req := &interaction.FavoriteListRequest{
		UserId: 12,
		Token:  c.Query("token"),
	}

	// 4. 发起RPC调用
	resp, err := cli.FavoriteList(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		panic(err)
	}

	// 5. gin返回给前端
	c.JSON(http.StatusOK, resp)
}

// @Summary xxx
// @Schemes
// @Description xxx
// @Tags 互动接口
// @Accept json
// @Produce json
// @Param request_body body relation.FavoriteActionRequest true "request body"
// @Router /douyin/comment/action/ [POST]

func InteractionCommentAction(c *gin.Context) {
	// 1. 创建客户端连接
	err := initInteractionCli()
	if err != nil {
		panic(err)
	}

	// 2. 创建发生消息的请求实例
	req := interaction.NewCommentActionRequest()

	// 3. 前端请求数据绑定到req中
	err = c.ShouldBindJSON(req)
	if err != nil {
		panic(err)
	}

	// 4. 发起RPC调用
	resp, err := cli.CommentAction(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		panic(err)
	}

	// 5. gin返回给前端
	c.JSON(http.StatusOK, resp)
}

// @Summary xxx
// @Schemes
// @Description xxx
// @Tags 互动接口
// @Accept json
// @Produce json
// @Param request_body body relation.FavoriteActionRequest true "request body"
// @Router /douyin/comment/list/ [GET]

func InteractionCommentList(c *gin.Context) {
	// 1. 创建客户端连接
	err := initInteractionCli()
	if err != nil {
		panic(err)
	}

	// 2. 创建发生消息的请求实例
	req := interaction.NewCommentListRequest()

	// 3. 前端请求数据绑定到req中
	err = c.ShouldBindJSON(req)
	if err != nil {
		panic(err)
	}

	// 4. 发起RPC调用
	resp, err := cli.CommentList(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		panic(err)
	}

	// 5. gin返回给前端
	c.JSON(http.StatusOK, resp)
}
