// Package video Package crud @author:戴林峰
// @date:2023/8/4
// @node:
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/crud"
	"github.com/douyin/common/oss"
	myRedis "github.com/douyin/common/redis"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
	"github.com/douyin/rpcServer/video/convert"
)

// 视频类型
const videoContentType = "application/mp4"

// 图片类型
const imageContentType = "image/png"

//	userPublishVideoList
//
// @Description: 根据用户id，获取缓存用户发布视频的key
// @param userId
// @return string
func userPublishVideoList(userId int) string {
	return fmt.Sprintf("video:feed:publish:%d", userId)
}

//	FindVideoListById
//
// @Description: 根据用户Id查找用户发布的视频
// 先在缓存当中查找，如果缓存没有就在数据库中查找，并更新缓存
// @param id
// @return []*models.Video
// @return error
func FindVideoListByUserId(userId int) ([]*model.Video, error) {
	cache, err := myRedis.NewRedisConn()
	if err != nil {
		log.Print("redis 客户端获取失败\n")
		LogCollector.Error(fmt.Sprintf("func user[%d]:FindVideoListByUserId Failed to get redis client in %s, err=%s", userId, time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return nil, err
	}
	errGet, _ := cache.Exists(context.Background(), userPublishVideoList(userId)).Result()
	//最终返回的video列表
	resVideoList := make([]*model.Video, 0)
	if errGet > 0 {
		// 缓存存在，直接从缓存中取出数据返回
		videoJsons, err := cache.LRange(context.Background(), userPublishVideoList(userId), 0, -1).Result()
		log.Printf("videoJsons.len = %+v\n", len(videoJsons))
		if err != nil {
			log.Printf("%v\n", err)
		}
		// 将json解析为video对象
		for _, videoJson := range videoJsons {
			videoDto := model.Video{}
			err := sonic.Unmarshal([]byte(videoJson), &videoDto)
			if err != nil {
				log.Fatalln("JSON decode error!")
				LogCollector.Error(fmt.Sprintf("func user[%d]:FindVideoListByUserId Failed to Unmarshal data  in %s, err=%s", userId, time.Now().Format("2006-01-02 15:04:05"), err.Error()))
			}
			resVideoList = append(resVideoList, &videoDto)
		}
	} else {
		//缓存不存在，查询数据库并写入缓存
		//从数据库中取出来的video列表
		videoList, err := models.FindVideoListBy("author_id", strconv.Itoa(userId))
		if err != nil {
			return nil, err
		}
		// 对kitex对象和model对象进行转换
		author, _ := crud.GetAuthor(uint(userId), uint(userId))
		resVideoList, err = convert.VideoSliceBo2Dto(videoList)
		for _, video := range resVideoList {
			video.Author = author
		}
		if err != nil {
			LogCollector.Error(fmt.Sprintf("func user[%d]:FindVideoListByUserId Failed to convert bo to dto  in %s, err=%s", userId, time.Now().Format("2006-01-02 15:04:05"), err.Error()))
			return nil, err
		}
		pipeline := cache.Pipeline()
		defer pipeline.Close()
		//依次对每一个视频对象进行序列化
		for _, video := range resVideoList {
			videoJson, _ := sonic.Marshal(video)
			_, err = pipeline.LPush(context.Background(), userPublishVideoList(userId), string(videoJson)).Result()
			if err != nil {
				LogCollector.Error(fmt.Sprintf("func user[%d]:FindVideoListByUserId Failed to execute redis cache  in %s, err=%s", userId, time.Now().Format("2006-01-02 15:04:05"), err.Error()))
				return nil, err
			}
		}
		pipeline.Expire(context.Background(), userPublishVideoList(userId), time.Minute)
		// 执行缓存操作
		pipeline.Exec(context.Background())
	}
	return resVideoList, nil
}

//	UploadVideo
//
// @Description: 执行视频上传逻辑
// 1. 上传视频到OSS
// 2. 将视频信息写入视频表
// 3. 更新用户发布视频的redis缓存，在list后面push当前的数据
// @param reader 视频文件的io流
// @param filename 文件名
// @param contentType 文件类型
// @param dataLen 数据长度
// @param userId 用户id
func UploadVideo(reader io.Reader, dataLen, userId int64, title string) error {
	service, _ := oss.GetOssService()
	snowId, _ := GenSnowId()
	//生成文件名
	filename := fmt.Sprintf("%s-%d", snowId, userId)
	// 视频文件名
	videoName := fmt.Sprintf("%s.mp4", filename)
	// 第一帧图片名
	imageName := fmt.Sprintf("%s.png", filename)
	//视频文件上传
	err := service.UploadFileWithBytestream(conf.MinioConfig.VideoBucketName, reader, videoName, dataLen, videoContentType)
	if err != nil {
		log.Fatalln("OSS视频文件上传失败")
		LogCollector.Error(fmt.Sprintf("func user[%d]:UploadVideo Failed to upload video to cos in %s, err=%s", userId, time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return err
	}
	// 获取到视频播放链接
	videoUrl, _ := service.GetPlayUrl(videoName)
	// 如果视频文件上传成功，那就获取视频文件url
	//抽取第一帧图片,得先上传视频文件，然后获取到视频文件的播放链接
	buffer, err := GetSnapshotImageBuffer(videoUrl, 1)
	if err != nil {
		log.Fatalln("视频封面图抽取失败！", err)
		LogCollector.Error(fmt.Sprintf("func user[%d]:UploadVideo Failed to get the first frame in %s, err=%s", userId, time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return err
	}
	//将第一帧图片转为io.reader
	imgReader := strings.NewReader(buffer.String())
	// 上传第一帧图片
	_ = service.UploadFileWithBytestream(conf.MinioConfig.VideoBucketName, imgReader, imageName, int64(buffer.Len()), imageContentType)
	// 获取第一帧图片的链接
	imageUrl, _ := service.GetPlayUrl(imageName)
	video := models.Video{
		Title:         title,
		FavoriteCount: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		AuthorID:      userId,
		PlayUrl:       videoUrl,
		CoverUrl:      imageUrl,
		CommentCount:  0,
	}
	//插入数据
	ID, err := models.InsertVideo(&video)
	bf.AddVideoId(ID)
	crud.DeletePublishListCache(int(userId))
	if err != nil {
		return err
	}
	// 将这一条数据插入redis
	cache, err := myRedis.NewRedisConn()
	if err != nil {
		log.Print("redis 客户端获取失败\n")
		return err
	}
	// 将video数据放入redis中，
	videoJson, _ := sonic.Marshal(video)
	//将这个视频添加到当前用户的发布视频cache当中去
	cache.RPush(context.Background(), userPublishVideoList(int(userId)), videoJson)
	return nil
}

//	GetVideoFeed
//
// @Description: 获取视频流
// @param latestTime
// @param nums
// @return []*models.Video
// @return error
// @return int64 返回上一次最后一个元素的时间
func GetVideoFeed(latestTime int64, nums int, userID uint) ([]*model.Video, error, int64) {
	cache, err := myRedis.NewRedisConn()
	if err != nil {
		log.Print("redis 客户端获取失败\n")
	}
	//缓存key
	cacheKey := fmt.Sprintf("video_feed_aa_%d", latestTime)
	cacheLastTimeKey := "video_feed_latest_time"
	errGet, err := cache.Exists(context.Background(), cacheKey).Result()
	// 最终返回的列表
	resVideoList := make([]*model.Video, 0)
	// 如果进入cache 这个flag就改为true，如果在cache执行的过程中有一个环节出错了，这个key就改为false。最后查询数据库
	cacheFlag := false
	var latestTimeRes int64 = 0
	// 表示缓存存在
	if errGet > 0 && latestTime != 0 {
		cacheFlag = true
		// 如果缓存存在，就直接从缓存中取数据返回
		videoJsons, err := cache.LRange(context.Background(), cacheKey, 0, -1).Result()
		log.Printf("videoJsons.len = %+v\n", len(videoJsons))
		if err != nil {
			log.Printf("%v\n", err)
			cacheFlag = false
		}
		// 将json解析为Video列表
		for _, videoJson := range videoJsons {
			videoDto := model.Video{}
			//将json字符串反序列化
			err := sonic.Unmarshal([]byte(videoJson), &videoDto)
			if err != nil {
				log.Printf("JSON decode error!")
				cacheFlag = false
			}
			resVideoList = append(resVideoList, &videoDto)
		}
		// 获取到当前对应列表的latest时间
		result, err := cache.Get(context.Background(), cacheLastTimeKey).Result()
		if err != nil {
			log.Print("获取当前feed列表对应的latest时间")
			cacheFlag = false
		}
		latestTimeRes, _ = strconv.ParseInt(result, 10, 64)
	}
	//如果走cache没有出错
	if cacheFlag {
		return resVideoList, nil, int64(latestTimeRes)
	}
	// 表示缓存不存在
	//获取latestTime时间之前的 不包括last？
	list, err := models.GetVideoFeedList(latestTime, nums)
	if err != nil {
		return nil, err, 0
	}
	resVideoList, err = convert.VideoSliceBo2Dto(list)
	if err != nil {
		return nil, err, 0
	}

	userVideoMap := map[int64]*model.User{}
	// crud, _ := crud.NewCachedCRUD()
	for i, item := range list {
		//如果这条视频已经查询过user了
		if v, ok := userVideoMap[item.ID]; ok {
			resVideoList[i].Author = v
			continue
		}
		//这一个还没被查询过
		// userBo, _ := models.GetUserById(item.AuthorID)
		// user, _ := convert.UserBo2Dto(*userBo)
		user, _ := crud.GetAuthor(userID, uint(item.AuthorID))
		isFavorite, _ := crud.IsFavorite(userID, uint(item.ID)) //TODO 没加缓存
		fmt.Println("isFavorite", isFavorite)
		//缓存
		userVideoMap[item.ID] = user
		resVideoList[i].Author = user
		resVideoList[i].IsFavorite = isFavorite
	}
	//当前列表的时间
	latestTimeRes = list[len(list)-1].CreatedAt.UnixMilli()

	// 到这一步的时候就需要将从mysql中查询出来的信息写入到redis中
	pipeline := cache.Pipeline()
	defer pipeline.Close()
	//依次对每一个视频对象进行序列化
	for _, video := range resVideoList {
		videoJson, _ := sonic.Marshal(video)
		_, err = pipeline.LPush(context.Background(), cacheKey, string(videoJson)).Result()
		if err != nil {
			log.Print("写缓存失败")
		}
	}
	latestTimeRes = 123456
	// 将当前播放列表的latestTime写入到cache中
	pipeline.Set(context.Background(), cacheLastTimeKey, latestTimeRes, 60)
	// 执行缓存操作
	pipeline.Exec(context.Background())
	return resVideoList, nil, latestTimeRes
}
