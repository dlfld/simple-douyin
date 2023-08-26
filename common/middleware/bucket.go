// Package middleware @author:戴林峰
// @date:2023/8/26
// @node:
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

// RateLimitMiddleware
//
//	@Description: Gin令牌桶限流
//	@param fillInterval 一秒钟
//	@param cap 容量
//	@param quantum 一秒钟填充多少个
//	@return gin.HandlerFunc
func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.String(http.StatusForbidden, "rate limit...")
			c.Abort()
			return
		}
		// 取道令牌就放行
		c.Next()
	}
}
