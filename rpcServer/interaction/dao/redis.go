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
			log.Printf("redis GetFavoriteVideoListByUserId key(%s) err(%+v)", key, err)
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

func (redis *redis) DelFavoriteVideoListByUserId(userId int64) error {
	key := fmt.Sprintf("FavoriteVideoList:userId:%d", userId)
	err := redis.cli.Del(context.Background(), key).Err()
	return err
}

func (redis *redis) GetCommentListByVideoId(videoId int64) (commentList []*model.Comment, err error) {
	key := fmt.Sprintf("FavoriteVideoList:videoId:%d", videoId)
	bs, err := redis.cli.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	} else {
		if err = json.Unmarshal(bs, &commentList); err != nil {
			log.Printf("redis GetCommentListByVideoId key(%s) err(%+v)", key, err)
			return nil, err
		}
	}
	return
}

func (redis *redis) SaveCommentListByVideoId(videoId int64, commentList []*model.Comment) error {
	key := fmt.Sprintf("FavoriteVideoList:videoId:%d", videoId)
	bs, _ := json.Marshal(commentList)
	err := redis.cli.Set(context.Background(), key, bs, ttl).Err()
	log.Printf("redis SaveCommentListByVideoId key(%s) err(%+v)", key, err)
	return err
}

func (redis *redis) DelCommentListByVideoId(videoId int64) error {
	key := fmt.Sprintf("FavoriteVideoList:videoId:%d", videoId)
	err := redis.cli.Del(context.Background(), key).Err()
	return err
}
