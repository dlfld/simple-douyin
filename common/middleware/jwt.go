package middleware

import (
	"fmt"

	"github.com/douyin/rpcServer/user/common"
	"github.com/gin-gonic/gin"
)

func JWT_AUTH(c *gin.Context) {
	Token := c.Query("token")
	if Token == "" {
		Token = c.PostForm("token") // 请求体中获取
	}
	common.ParseToken(Token)
	_, claims, err1 := common.ParseToken(Token)
	if err1 != nil {
		panic(err1)
	}
	useridClaims := claims.UserId
	fmt.Println("jwt:", useridClaims)
	c.Set("userID", useridClaims)
	c.Next()
}
