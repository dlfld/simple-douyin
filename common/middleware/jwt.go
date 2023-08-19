package middleware

import (
	"fmt"
	"github.com/douyin/rpcServer/user/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWT_AUTH(c *gin.Context) {
	Token := c.Query("token")
	if Token == "" {
		tokenMap := make(map[string]string, 0)
		err := c.BindJSON(&tokenMap)
		if err != nil {
			c.JSON(http.StatusOK, noToken)
			c.Abort()
			return
		}
		Token = tokenMap["token"]
	}
	_, claims, err1 := common.ParseToken(Token)
	if err1 != nil {
		c.JSON(http.StatusOK, invalidToken)
		c.Abort()
	}
	useridClaims := claims.UserId
	fmt.Println("jwt:", useridClaims)
	c.Set("userID", useridClaims)
	c.Next()
}

type resp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

var noToken = resp{1, "缺少token"}
var invalidToken = resp{1, "无效token"}
