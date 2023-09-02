package crud

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
	"github.com/go-redis/redis/v8"
)

func userWatchedVideo(userID int64) string {
	return fmt.Sprintf("video:feed:watched:%d", userID)
}
func hotVideo() string {
	return "video:feed:hot"
}
func videosCache(VideoID int64) string {
	return fmt.Sprintf("video:cache:%d", VideoID)
}

// 缓存用户看过哪些视频 防止出现重复结果
func CacheUserWatched(userID int64, videoIDs []int64) (err error) {
	pipline := crud.redis.Pipeline()
	// 设置过期时间

	for _, v := range videoIDs {
		pipline.SAdd(context.Background(), userWatchedVideo(userID), v)
	}
	pipline.Expire(context.Background(), userWatchedVideo(userID), time.Hour)
	_, err = pipline.Exec(context.Background())
	return
}

// 如果所有视频都看过了就删除缓存 防止没有数据返回
func DeleteUserWatched(userID int64) (err error) {
	err = crud.redis.Del(context.Background(), userWatchedVideo(userID)).Err()
	return
}

// GetUserUnWatched 获取用户没有看过的videoID
func GetUserUnWatched(userID int64) (videoIDs []int, err error) {

	// 加载最新的n个视频id到集合 视频推广也可放入此缓存中 会优先推送
	if crud.redis.Exists(context.Background(), hotVideo()).Val() != 1 {
		var vids []any
		err = crud.mysql.Model(&models.Video{}).Select("id").Limit(500).Find(&vids).Error
		if err != nil {
			return
		}

		err = crud.redis.SAdd(context.Background(), hotVideo(), vids...).Err()
		// 设置一天过期时间
		crud.redis.Expire(context.Background(), hotVideo(), time.Hour*24)
		if err != nil {
			return
		}
	}
	// 差集计算没有看过的视频
	idStrList := crud.redis.SDiff(context.Background(), hotVideo(), userWatchedVideo(userID)).Val()
	n := len(idStrList)
	videoIDs = make([]int, n)
	// Set是按照升序排列的  倒序转换为降序
	for i, v := range idStrList {
		videoIDs[n-i-1], _ = strconv.Atoi(v)
	}
	return
}

func CacheVideos(videos []*models.Video) (err error) {
	pipline := crud.redis.Pipeline()
	for _, v := range videos {
		pipline.HSet(context.Background(), videosCache(v.ID),
			"id", v.ID,
			"author_id", v.AuthorID,
			"play_url", v.PlayUrl,
			"cover_url", v.CoverUrl,
			"title", v.Title,
			"favorite_count", v.FavoriteCount,
			"comment_count", v.CommentCount)
		pipline.Expire(context.Background(), videosCache(v.ID), time.Hour*24)

	}
	_, err = pipline.Exec(context.Background())
	return
}

func GetVideos(videoIDs []int) (videos []*models.Video, err error) {
	videos = make([]*models.Video, len(videoIDs))
	pipline := crud.redis.Pipeline()
	for _, v := range videoIDs {
		pipline.HGetAll(context.Background(), videosCache(int64(v)))
	}
	results, err := pipline.Exec(context.Background())
	if err != nil {
		return
	}
	uncached := make([]int, 0)
	uncached_pos := make([]int, 0)
	// TODO fix bug
	for i, v := range results {

		if v.Err() != nil {
			uncached = append(uncached, videoIDs[i])
			uncached_pos = append(uncached_pos, i)
			continue
		}
		videos[i] = &models.Video{}
		err := v.(*redis.StringStringMapCmd).Scan(videos[i])
		resultMap, err := v.(*redis.StringStringMapCmd).Result()
		if len(resultMap) == 0 || err != nil {
			uncached = append(uncached, videoIDs[i])
			uncached_pos = append(uncached_pos, i)
			continue
		}
	}
	if len(uncached) > 0 {
		var DBModels []*models.Video
		crud.mysql.Model(&models.Video{}).Where("id in (?)", uncached).Find(&DBModels)
		for i, v := range uncached_pos {
			videos[v] = DBModels[i]
		}
		CacheVideos(DBModels)
	}

	return
}

// GetUserFeed TODO 缓存Video对象 避免反复查询数据库
// 返回当前用户的视频流
func GetUserFeed(UserID int64, latestTime int64) (videos []*model.Video, next_time int64, err error) {
	// 查询当前用户没有看过的热点视频记录
	unwatched, err := GetUserUnWatched(UserID)
	if err != nil {
		unwatched = make([]int, 0)
		err = nil
	}
	var DBModels []*models.Video
	// 查询没有看过的热点视频记录
	if len(unwatched) > 0 {
		// crud.mysql.Model(&models.Video{}).Order("id desc").Where("id in (?)", unwatched).Find(&DBModels)
		DBModels, err = GetVideos(unwatched)
		if err != nil {
			return
		}
	} else { // 数量=0
		// 查询数据库
		crud.mysql.Raw("SELECT * FROM videos WHERE created_at < ? ORDER BY created_at DESC LIMIT ?;", time.UnixMilli(latestTime), 30).Find(&DBModels)
		CacheVideos(DBModels)
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
	isFavs, err := IsFavorites(UserID, videoIDs)
	if err != nil {
		isFavs = make(map[int64]bool, 0)
		err = nil
	}
	//fmt.Printf("UserId = %d,autherIDS = %d", UserID, autherIDs)

	authors, _ := GetAuthors(UserID, autherIDs)
	for i, v := range DBModels {
		videos[i].Author = authors[v.AuthorID]
		videos[i].IsFavorite = isFavs[v.ID]
		// videos[i].IsFavorite, _ = IsFavorite(uint(UserID), uint(v.ID))
	}
	if len(DBModels) > 1 {
		next_time = DBModels[len(DBModels)-1].CreatedAt.UnixMilli()
	}
	// 记录用户看过的视频
	err = CacheUserWatched(UserID, videoIDs)
	if err != nil {
		return
	}
	return
}
