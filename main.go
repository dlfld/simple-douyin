package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/douyin/common/middleware"
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
	//对桶进行设置
	bucket := middleware.RateLimitMiddleware(1*time.Second, 10000, 10)
	apiRouter := r.Group("/douyin")
	jwtRouter := r.Group("/douyin")
	limitTestRouter := r.Group("/douyin")

	jwtRouter.Use(middleware.JWT_AUTH, bucket)
	limitTestRouter.Use(bucket)
	// 添加中间件
	//jwtRouter.Use(middleware.MaxAllowed)
	// 只用来解析token不做拦截
	apiRouter.Use(bucket)
	apiRouter.Use(middleware.JWT_PARSE)

	apiRouter.GET("/user/", handler.UserInfo)
	apiRouter.POST("/user/register/", handler.Register)
	apiRouter.POST("/user/login/", handler.Login)

	jwtRouter.POST("/relation/action/", handler.RelationAction)
	apiRouter.GET("/relation/follow/list/", handler.RelationFollowList)
	apiRouter.GET("/relation/follower/list/", handler.RelationFollowerList)
	apiRouter.GET("/relation/friend/list/", handler.RelationFriendList)

	jwtRouter.GET("/message/chat/", handler.MessageChat)
	jwtRouter.POST("/message/action/", handler.MessageAction)

	// 视频相关接口
	apiRouter.GET("/publish/list/", handler.PublishList)
	jwtRouter.POST("/publish/action/", handler.VideoSubmit)
	apiRouter.GET("/feed/", handler.VideoFeed)

	//互动interaction
	jwtRouter.POST("/favorite/action/", handler.InteractionFavoriteAction)
	apiRouter.GET("/favorite/list/", handler.InteractionFavoriteList)
	jwtRouter.POST("/comment/action/", handler.InteractionCommentAction)
	apiRouter.GET("/comment/list/", handler.InteractionCommentList)

	limitTestRouter.GET("/test", handler.BucketLimit)

	handler.InitRpcCli()
	r.Run("0.0.0.0:8080")
}
