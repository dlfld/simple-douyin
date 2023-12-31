package crud

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/douyin/common/gorse"
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

func DeletePublishListCache(userID int) (err error) {
	err = crud.redis.Del(context.Background(), fmt.Sprintf("video:feed:publish:%d", userID)).Err()
	return
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
func GetUserUnWatched(userID, limit int64) (videoIDs []int, err error) {

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
	if n > int(limit) {
		videoIDs = make([]int, limit)
		for i, v := range idStrList[n-int(limit):] {
			// Set是按照升序排列的  倒序转换为降序
			videoIDs[limit-1-int64(i)], _ = strconv.Atoi(v)
		}
		return
	} else {
		videoIDs = make([]int, n)
		for i, v := range idStrList {
			// Set是按照升序排列的  倒序转换为降序
			videoIDs[n-1-i], _ = strconv.Atoi(v)
		}
		return
	}

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
			"comment_count", v.CommentCount,
			"created_at", v.CreatedAt)
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
		videos[i] = new(models.Video)
		m := v.(*redis.StringStringMapCmd).Val()
		if len(m) == 0 {
			uncached = append(uncached, videoIDs[i])
			uncached_pos = append(uncached_pos, i)
			continue
		}
		b, err := sonic.Marshal(m)
		if err != nil {
			uncached = append(uncached, videoIDs[i])
			uncached_pos = append(uncached_pos, i)
			continue
		}
		err = sonic.Unmarshal(b, &videos[i])
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

func GetRecommend(UserID, limit int) (videos []*model.Video, next_time int64, err error) {
	// 从gorse获取推荐结果 自动设置已读
	videoIdStrList, err := gorse.Client.GetItemRecommend(context.Background(), strconv.Itoa(int(UserID)), nil, "read", "1m", limit, 0)
	if err != nil {
		return
	}
	// 转换为int数组

	videoIds := make([]int, len(videoIdStrList))
	for i, v := range videoIdStrList {
		n, e := strconv.Atoi(v)
		if e != nil {
			continue
		}
		videoIds[i] = n
	}
	// 查询视频详情
	DBModels, err := GetVideos(videoIds)
	if err != nil {
		return
	}
	videoIDs := make([]int64, len(DBModels))
	autherIDs := make([]int64, len(DBModels))
	videos = make([]*model.Video, len(DBModels))
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
	// 查询是否被当前用户关注
	isFavs, err := IsFavorites(int64(UserID), videoIDs)
	if err != nil {
		isFavs = make(map[int64]bool, 0)
		err = nil
	}

	// 获取视频作者信息
	authors, _ := GetAuthors(int64(UserID), autherIDs)
	for i, v := range DBModels {
		videos[i].Author = authors[v.AuthorID]
		videos[i].IsFavorite = isFavs[v.ID]

	}
	// 时间戳
	if len(DBModels) > 1 {
		next_time = DBModels[len(DBModels)-1].CreatedAt.UnixMilli()
	}
	return
}

// GetUserFeed
// 返回当前用户的视频流
func GetUserFeed(UserID int64, latestTime int64) (videos []*model.Video, next_time int64, err error) {
	videosQueryed := 0
	// 使用gorse查询推荐部分视频

	videos, next_time, _ = GetRecommend(int(UserID), 20)
	// 第一个视频直接设置为已读
	if len(videos) > 0 {
		gorse.Client.InsertFeedback(context.Background(), []gorse.Feedback{{
			FeedbackType: "read",
			UserId:       strconv.Itoa(int(UserID)),
			ItemId:       strconv.Itoa(int(videos[0].Id)),
			Timestamp:    time.Now().Format("2006-01-02 15:04:05")}})
	} else {
		videos = make([]*model.Video, 0)
	}
	// 防止出现重复视频
	recommendIDs := make([]int64, len(videos))
	for i, v := range videos {
		recommendIDs[i] = v.Id
	}
	CacheUserWatched(UserID, recommendIDs)
	videosQueryed += len(videos)
	// 查询当前用户没有看过的活动视频记录
	unwatched, err := GetUserUnWatched(UserID, 30-int64(videosQueryed))
	if err != nil {
		unwatched = make([]int, 0)
		err = nil
	}
	var DBModels []*models.Video
	// 查询没有看过的活动视频记录
	if len(unwatched) > 0 {
		// crud.mysql.Model(&models.Video{}).Order("id desc").Where("id in (?)", unwatched).Find(&DBModels)
		DBModels, err = GetVideos(unwatched)
		videosQueryed += len(DBModels)
		if err != nil {
			return
		}
	}
	for _, v := range DBModels {
		recommendIDs = append(recommendIDs, v.ID)
	}
	if videosQueryed < 30 { // 数量<30 从数据库中查询补充
		// 查询数据库
		var MysqlResults []*models.Video
		if latestTime == 0 {
			if len(recommendIDs) == 0 {
				crud.mysql.Model(&models.Video{}).Order("id DESC").Limit(30 - videosQueryed).Find(&MysqlResults)
			} else { // 避免重复推荐
				crud.mysql.Model(&models.Video{}).Where("id NOT IN (?)", recommendIDs).Order("id DESC").Limit(30 - videosQueryed).Find(&MysqlResults)
			}
		} else {
			if len(recommendIDs) == 0 {
				crud.mysql.Raw("SELECT * FROM videos WHERE created_at < ?  ORDER BY created_at DESC LIMIT ?;", time.UnixMilli(latestTime), 30-videosQueryed).Find(&MysqlResults)
				if len(MysqlResults) == 0 {
					crud.mysql.Model(&models.Video{}).Order("id DESC").Limit(30 - videosQueryed).Find(&MysqlResults)
				}
			} else { // 避免重复推荐
				crud.mysql.Raw("SELECT * FROM videos WHERE created_at < ? AND id NOT IN ? ORDER BY created_at DESC LIMIT ?;", time.UnixMilli(latestTime), recommendIDs, 30-videosQueryed).Find(&MysqlResults)
				if len(MysqlResults) == 0 {
					crud.mysql.Model(&models.Video{}).Where("id NOT IN (?)", recommendIDs).Order("id DESC").Limit(30 - videosQueryed).Find(&MysqlResults)
				}
			}

		}
		CacheVideos(MysqlResults)
		DBModels = append(DBModels, MysqlResults...)
	}

	videoIDs := make([]int64, len(DBModels))
	autherIDs := make([]int64, len(DBModels))
	videosAdv := make([]*model.Video, len(DBModels))
	// 遍历已观看视频用户ID
	for i, v := range DBModels {
		videosAdv[i] = &model.Video{
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

	authors, _ := GetAuthors(UserID, autherIDs)

	for i, v := range DBModels {
		videosAdv[i].Author = authors[v.AuthorID]
		videosAdv[i].IsFavorite = isFavs[v.ID]
	}

	// 记录用户看过的视频
	err = CacheUserWatched(UserID, videoIDs)
	// 合并推荐视频和活动视频
	videos = append(videos, videosAdv...)
	if len(videos) > 0 {
		res, e := GetVideos([]int{int(videos[len(videos)-1].Id)})
		if e == nil && len(res) > 0 {
			next_time = res[0].CreatedAt.UnixMilli()
		}
	}

	if err != nil {
		return
	}
	return
}
