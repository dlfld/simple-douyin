package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestGetLatestTime(t *testing.T) {
	key := fmt.Sprintf("%d:%d", 1111, 22)
	res, err := cache.getMinScore(context.Background(), key)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(res)
}

func TestCache_ZRangeByScore(t *testing.T) {
	key := fmt.Sprintf("%d:%d", 1111, 22)
	res := cache.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
		Min: "0",
		Max: "inf",
	})
	if cache.Err() != nil {
		t.Error(cache.Err().Error())
	}
	for i := 0; i < len(res); i++ {
		fmt.Println(res[i])
	}
}

func TestCacheGetMimCache(t *testing.T) {
	key := fmt.Sprintf("%s:%d:%d", messageCacheTable, 1111, 22)
	var res int64 = 0
	cache.getRecordNum(context.Background(), key, &res)
	if cache.Err() != nil {
		t.Error(cache.Err().Error())
	}
	fmt.Println(res)

}
