// Package handler @author:戴林峰
// @date:2023/8/26
// @node:
package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BucketLimit(c *gin.Context) {
	c.JSON(http.StatusOK, "访问到了。。。")
}
