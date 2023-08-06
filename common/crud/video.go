// Package crud @author:戴林峰
// @date:2023/8/4
// @node:
package crud

import (
	"context"
	"fmt"
	myRedis "github.com/douyin/common/redis"
	"github.com/douyin/models"
	"github.com/go-redis/redis/v8"
	"log"
)

//	userPublishVideoList
//
// @Description: 根据用户id，获取缓存用户发布视频的key
// @param userId
// @return string
func userPublishVideoList(userId int) string {
	return fmt.Sprintf("user_video_publish_%d", userId)
}

//	FindVideoListById
//
// @Description: 根据用户Id查找用户发布的视频
// 先在缓存当中查找，如果缓存没有就在数据库中查找，并更新缓存
// @param id
// @return []*models.Video
// @return error
func FindVideoListByUserId(userId int) ([]*models.Video, error) {
	cache, err := myRedis.NewRedisConn()
	if err != nil {
		log.Print("redis 客户端获取失败\n")
		return nil, err
	}
	_, errGet := cache.Get(context.Background(), userPublishVideoList(userId)).Result()

	if errGet == redis.Nil {
		//缓存不存在

	} else if errGet != nil {
		return nil, errGet
	}
	// 缓存存在
	return nil, nil
}
