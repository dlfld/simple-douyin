package crud

import (
	"context"
	"fmt"
	"log"

	"github.com/douyin/models"
	"github.com/go-redis/redis/v8"
)

func favoriteVideos(userID int64) string {
	return fmt.Sprintf("interaction:favoriteVideos:%d", userID)
}

func FavoriteVideo(userID, videoID int64) (err error) {
	if crud.redis.Exists(context.Background(), favoriteVideos(userID)).Val() != 1 {
		var vids []any
		err = crud.mysql.Model(&models.FavoriteVideoRelation{}).Select("video_id").Where("user_id=?", userID).Find(&vids).Error
		if err != nil {
			return
		}
		err = crud.redis.SAdd(context.Background(), favoriteVideos(userID), vids...).Err()
		if err != nil {
			return
		}
	}
	err = crud.redis.SAdd(context.Background(), favoriteVideos(userID), videoID).Err()
	return
}

func UnFavoriteVideo(userID, videoID int64) (err error) {
	if crud.redis.Exists(context.Background(), favoriteVideos(userID)).Val() != 1 {
		var vids []any
		err = crud.mysql.Model(&models.FavoriteVideoRelation{}).Select("video_id").Where("user_id=?", userID).Find(&vids).Error
		if err != nil {
			return
		}
		err = crud.redis.SAdd(context.Background(), favoriteVideos(userID), vids...).Err()
		if err != nil {
			return
		}
	}
	err = crud.redis.SRem(context.Background(), favoriteVideos(userID), videoID).Err()
	return
}

// IsFavorite 判断是否点赞
func IsFavorite(self uint, videoId uint) (isFavorite bool, err error) {
	var i int64
	result := crud.mysql.Raw("select 1 from user_favorite_videos WHERE user_id = ? AND video_id = ? LIMIT 1", self, videoId).Count(&i)

	if result.Error != nil {
		log.Println(err)
	}
	return i != 0, result.Error
}

// // IsFavorite 判断是否点赞
func IsFavorites(self int64, videoId []int64) (videoFav map[int64]bool, err error) {
	// var i int64
	if crud.redis.Exists(context.Background(), favoriteVideos(self)).Val() != 1 {
		var vids []any
		err = crud.mysql.Model(&models.FavoriteVideoRelation{}).Select("video_id").Where("user_id=?", self).Find(&vids).Error
		if err != nil {
			return
		}
		err = crud.redis.SAdd(context.Background(), favoriteVideos(self), vids...).Err()
		if err != nil {
			return
		}
	}
	pipline := crud.redis.Pipeline()
	for _, v := range videoId {
		pipline.SIsMember(context.Background(), favoriteVideos(self), v)
	}
	result, err := pipline.Exec(context.Background())
	if err != nil {
		return
	}
	videoFav = make(map[int64]bool)
	for i, v := range result {
		videoFav[videoId[i]] = v.(*redis.BoolCmd).Val()
	}
	return videoFav, nil
}
