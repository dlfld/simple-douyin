package dao

import (
	"context"
	"encoding/json"
	"fmt"
	comrd "github.com/douyin/common/redis"
	"github.com/douyin/kitex_gen/model"
	"log"
	"time"

	grd "github.com/go-redis/redis/v8"
)

const ttl = 30 * time.Minute

type redis struct {
	cli *grd.Client
}

func NewRedis() *redis {
	conn, _ := comrd.NewRedisConn()
	return &redis{cli: conn}
}

func (redis *redis) Test() {
	_ = redis.cli.String()
}

func (redis *redis) GetFavoriteVideoListByUserId(userId int64) (videoList []*model.Video, err error) {
	key := fmt.Sprintf("FavoriteVideoList:userId:%d", userId)
	bs, err := redis.cli.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	} else {
		if err = json.Unmarshal(bs, &videoList); err != nil {
			return nil, err
		}
	}
	return
}

func (redis *redis) SaveFavoriteVideoListByUserId(userId int64, videoList []*model.Video) error {
	key := fmt.Sprintf("FavoriteVideoList:userId:%d", userId)
	bs, _ := json.Marshal(videoList)
	err := redis.cli.Set(context.Background(), key, bs, ttl).Err()
	log.Printf("redis SaveFavoriteVideoListByUserId key(%s) err(%+v)", key, err)
	return err
}
