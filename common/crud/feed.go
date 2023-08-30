package crud

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
	"github.com/go-redis/redis/v8"
)

func userWatchedVideo(userID int64) string {
	return fmt.Sprintf("video:feed:watched:%d", userID)
}
func totalVideo() string {
	return "video:feed:total"
}
func videosCache() string {
	return "video:feed:videos"
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

// GetUserUnWatched 获取用户没有看过的videoID
func GetUserUnWatched(userID int64) (videoIDs []int, err error) {
	// 加载所有视频id到集合
	if crud.redis.Exists(context.Background(), totalVideo()).Val() != 1 {
		var vids []any
		err = crud.mysql.Model(&models.Video{}).Select("id").Find(&vids).Error
		if err != nil {
			return
		}
		err = crud.redis.SAdd(context.Background(), totalVideo(), vids...).Err()
		if err != nil {
			return
		}
	}
	// 差集计算没有看过的视频
	idStrList := crud.redis.SDiff(context.Background(), totalVideo(), userWatchedVideo(userID)).Val()
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
		data, err := sonic.Marshal(v)
		if err != nil {
			return err
		}
		pipline.HSet(context.Background(), videosCache(), v.ID, string(data))
	}
	_, err = pipline.Exec(context.Background())
	return
}

func GetVideos(videoIDs []int) (videos []*models.Video, err error) {
	videos = make([]*models.Video, len(videoIDs))
	pipline := crud.redis.Pipeline()
	for _, v := range videoIDs {
		pipline.HGet(context.Background(), videosCache(), strconv.Itoa(int(v)))
	}
	results, err := pipline.Exec(context.Background())
	if err != nil {
		return
	}
	uncached := make([]int, 0)
	uncached_pos := make([]int, 0)
	for i, v := range results {
		if v.Err() != nil {
			uncached = append(uncached, videoIDs[i])
			uncached_pos = append(uncached_pos, i)
			continue
		}
		data, err := v.(*redis.StringCmd).Result()
		if err != nil {
			uncached = append(uncached, videoIDs[i])
			uncached_pos = append(uncached_pos, i)
			continue
		}
		videos[i] = &models.Video{}
		err = sonic.Unmarshal([]byte(data), &videos[i])
		if err != nil {
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
func GetUserFeed(UserID int64, latestTime int64) (videos []*model.Video, err error) {
	// 查询当前用户没有看过的视频记录
	unwatched, _ := GetUserUnWatched(UserID)
	var DBModels []*models.Video
	// 查询没有看过的视频记录
	if len(unwatched) >= 30 {
		// crud.mysql.Model(&models.Video{}).Order("id desc").Where("id in (?)", unwatched[:30]).Find(&DBModels)
		DBModels, err = GetVideos(unwatched[:30])
		if err != nil {
			return
		}
	} else if len(unwatched) == 0 { // 数量=0
		DeleteUserWatched(UserID)
		crud.mysql.Model(&models.Video{}).Order("id desc").Limit(30).Find(&DBModels)
	} else {
		crud.mysql.Model(&models.Video{}).Order("id desc").Where("id in (?)", unwatched).Find(&DBModels)
		// 数量小于30 首尾相连 避免重复出现同一个数据
		var appendDB []*models.Video
		crud.mysql.Model(&models.Video{}).Order("id desc").Where("id not in (?)", unwatched).Limit(30 - len(unwatched)).Find(&appendDB)
		DBModels = append(DBModels, appendDB...)
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
		return
	}
	//fmt.Printf("UserId = %d,autherIDS = %d", UserID, autherIDs)

	authors, _ := GetAuthors(UserID, autherIDs)
	for i, v := range DBModels {
		videos[i].Author = authors[v.AuthorID]
		videos[i].IsFavorite = isFavs[v.ID]
		// videos[i].IsFavorite, _ = IsFavorite(uint(UserID), uint(v.ID))
	}
	// 记录用户看过的视频
	err = CacheUserWatched(UserID, videoIDs)
	if err != nil {
		return
	}
	return
}
