package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/douyin/common/middleware"
	"github.com/douyin/controller"

	docs "github.com/douyin/docs"
	"github.com/douyin/handler"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func getWorkingDirPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

func genSwagger() {
	// 有问题，不能用
	sysType := runtime.GOOS
	if sysType == "linux" || sysType == "darwin" {
		absWd := getWorkingDirPath()
		command := "bash " + absWd + "/bash/swag_gen.sh"
		cmd := exec.Command(command)
		err := cmd.Run()
		fmt.Printf("%+v\n", err)
	}
}
func main() {
	//执行生成swagger文件的命令 warning 失效
	//genSwagger()
	// public directory is used to serve static resources
	//r.Static("/static", "./public")
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	apiRouter := r.Group("/douyin")
	jwtRouter := r.Group("/douyin")
	jwtRouter.Use(middleware.JWT_AUTH)
	// 只用来解析token不做拦截
	apiRouter.Use(middleware.JWT_PARSE)
	//// basic apis
	//apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	//apiRouter.POST("/publish/action/", controller.Publish)
	//apiRouter.GET("/publish/list/", controller.PublishList)
	//
	//// extra apis - I
	//apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	//apiRouter.GET("/favorite/list/", controller.FavoriteList)
	//apiRouter.POST("/comment/action/", controller.CommentAction)
	//apiRouter.GET("/comment/list/", controller.CommentList)
	//
	//// extra apis - II
	jwtRouter.POST("/relation/action/", handler.RelationAction)
	apiRouter.GET("/relation/follow/list/", handler.RelationFollowList)
	apiRouter.GET("/relation/follower/list/", handler.RelationFollowerList)
	apiRouter.GET("/relation/friend/list/", handler.RelationFriendList)

	jwtRouter.GET("/message/chat/", handler.MessageChat)
	jwtRouter.POST("/message/action/", handler.MessageAction)

	// 视频相关接口
	jwtRouter.GET("/publish/list/", handler.PublishList)
	jwtRouter.POST("/publish/action/", handler.VideoSubmit)
	apiRouter.GET("/feed/", handler.VideoFeed)
	// apiRouter.GET("/t/ ", handler.RelationFollowerList)

	//互动interaction
	jwtRouter.POST("/favorite/action/", handler.InteractionFavoriteAction)
	apiRouter.GET("/favorite/list/", handler.InteractionFavoriteList)
	jwtRouter.POST("/comment/action/", handler.InteractionCommentAction)
	apiRouter.GET("/comment/list/", handler.InteractionCommentList)

	err := handler.InitInteractionCli()
	if err != nil {
		return
	}
	r.Run("0.0.0.0:8080")
}
