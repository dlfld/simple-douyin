package middleware

import (
	"net/http"

	"github.com/douyin/rpcServer/user/common"
	"github.com/gin-gonic/gin"
)

func JWT_AUTH(c *gin.Context) {
	Token, ok := c.GetQuery("token")
	if !ok {
		Token = c.PostForm("token")
	}
	if Token == "" && c.Request.Method == "POST" {
		tokenMap := make(map[string]any, 0)
		err := c.ShouldBindJSON(&tokenMap)
		if err != nil {
			c.JSON(http.StatusForbidden, noToken)
			c.Abort()
			return
		}
		Token = tokenMap["token"].(string)
	}
	_, claims, err1 := common.ParseToken(Token)
	if err1 != nil {
		c.JSON(http.StatusForbidden, invalidToken)
		c.Abort()
	}
	useridClaims := claims.UserId
	c.Set("userID", useridClaims)
	c.Next()
}

func JWT_PARSE(c *gin.Context) {
	Token, ok := c.GetQuery("token")
	if !ok {
		Token = c.PostForm("token")
	}
	if Token == "" && c.Request.Method == "POST" {
		tokenMap := make(map[string]string, 0)
		err := c.ShouldBindJSON(&tokenMap)
		if err != nil {
			return
		}
		Token = tokenMap["token"]
	}
	_, claims, err1 := common.ParseToken(Token)
	if err1 != nil {
		return
	}
	useridClaims := claims.UserId
	c.Set("userID", useridClaims)
}

type resp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

var noToken = resp{1, "缺少token"}
var invalidToken = resp{1, "无效token"}
