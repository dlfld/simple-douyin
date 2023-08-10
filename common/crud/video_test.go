// @author:戴林峰
// @date:2023/8/5
// @node:
package crud

import (
	"context"
	"fmt"
	myRedis "github.com/douyin/common/redis"
	"testing"
)

func TestFindVideoListByUserId(t *testing.T) {
	cache, _ := myRedis.NewRedisConn()
	exists := cache.Exists(context.Background(), "aaa").
		// print("aaa\n")
	result, _ := cache.Get(context.Background(), "aaa").Result()
	fmt.Printf("result:%+v\n", result)
	s2 := cache.LRange(context.Background(), "mylist", 0, -1).Val()
	fmt.Printf("result:%+v\n", s2)
}
