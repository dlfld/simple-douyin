package crud

import (
	"context"
	"fmt"
	"time"

	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
)

func userWatchedVideo(userID int64) string {
	return fmt.Sprintf("video:feed:watched:%d", userID)
}

// 缓存用户看过哪些视频 防止出现重复结果
func CacheUserWatched(userID int64, videoIDs []int64) (err error) {
	pipline := crud.redis.Pipeline()
	// 设置过期时间

	for _, v := range videoIDs {
		pipline.SAdd(context.Background(), userWatchedVideo(userID), v)
	}
	pipline.ExpireXX(context.Background(), userWatchedVideo(userID), time.Hour)
	_, err = pipline.Exec(context.Background())
	return
}

// 如果所有视频都看过了就删除缓存 防止没有数据返回
func DeleteUserWatched(userID int64) (err error) {
	err = crud.redis.Del(context.Background(), userWatchedVideo(userID)).Err()
	return
}

// 获取已经看过的videoID
func GetUserWatched(userID int64) (videoIDs []string, err error) {
	videoIDs = crud.redis.SMembers(context.Background(), userWatchedVideo(userID)).Val()
	return
}

// TODO 新发布的视频会被重复刷到两次
// 返回当前用户的视频流
func GetUserFeed(UserID int64, latestTime int64) (videos []*model.Video, err error) {
	// 查询当前用户看过那些视频
	watched, _ := GetUserWatched(UserID)
	var DBModels []*models.Video
	// 查询没有看过的视频记录
	if len(watched) > 0 {
		var cnt int64
		crud.mysql.Model(&models.Video{}).Order("id desc").Where("id not in (?)", watched).Limit(30).Count(&cnt)
		if cnt < 30 { // 没有匹配的视频 就删除缓存重新查询
			DeleteUserWatched(UserID)
			crud.mysql.Model(&models.Video{}).Order("id desc").Limit(30).Find(&DBModels)
		} else {
			crud.mysql.Model(&models.Video{}).Order("id desc").Where("id not in (?)", watched).Limit(30).Find(&DBModels)
		}
	} else { // watched为空就不用not in条件
		crud.mysql.Model(&models.Video{}).Order("id desc").Limit(30).Find(&DBModels)
	}

	videoIDs := make([]int64, len(DBModels))
	autherIDs := make([]int64, len(DBModels))
	videos = make([]*model.Video, len(DBModels))
	// 遍历已观看视频用户ID
	for i, v := range DBModels {
		videos[i] = &model.Video{
			Id:            v.ID,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			Title:         v.Title,
		}
		videoIDs[i] = v.ID
		autherIDs[i] = v.AuthorID
	}
	//TODO 批量处理isFavourite
	// isFavs=GetFavourites(userID,videoIDs)
	authors, _ := GetAuthors(UserID, autherIDs)
	for i, v := range DBModels {
		videos[i].Author = authors[v.AuthorID]
		// videos[i].IsFavorite=isFavs[v.id]
		videos[i].IsFavorite, _ = IsFavorite(uint(UserID), uint(v.ID))
	}
	// 记录用户看过的视频
	err = CacheUserWatched(UserID, videoIDs)
	if err != nil {
		return
	}
	return
}
