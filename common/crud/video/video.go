// Package video Package crud @author:戴林峰
// @date:2023/8/4
// @node:
package video

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/crud/video/convert"
	"github.com/douyin/common/oss"
	myRedis "github.com/douyin/common/redis"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
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
	return fmt.Sprintf("user_video_publish_%d", userId)
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
		return nil, err
	}
	errGet, err := cache.Exists(context.Background(), userPublishVideoList(userId)).Result()
	//最终返回的video列表
	//var resVideoList []*model.Video = nil
	resVideoList := make([]*model.Video, 0)
	if errGet > 0 {
		// 缓存存在，直接从缓存中取出数据返回
		videoJsons, err := cache.LRange(context.Background(), userPublishVideoList(userId), 0, -1).Result()
		log.Printf("videoJsons.len = %+v\n", len(videoJsons))
		if err != nil {
			panic(err)
		}
		// 将json解析为video对象
		for _, videoJson := range videoJsons {
			videoDto := model.Video{}
			err := json.Unmarshal([]byte(videoJson), &videoDto)
			if err != nil {
				log.Fatalln("JSON decode error!")
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
		resVideoList, err = convert.VideoSliceBo2Dto(videoList)
		if err != nil {
			return nil, err
		}
		pipeline := cache.Pipeline()
		defer pipeline.Close()
		//依次对每一个视频对象进行序列化
		for _, video := range resVideoList {
			videoJson, _ := json.Marshal(video)
			_, err = pipeline.LPush(context.Background(), userPublishVideoList(userId), string(videoJson)).Result()
			if err != nil {
				return nil, err
			}
		}
		// 执行缓存操作
		_, err = pipeline.Exec(context.Background())
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
// @param videoUrl 视频url
// @param dataLen 数据长度
// @param userId 用户id
func UploadVideo(reader io.Reader, filename, videoUrl string, dataLen, userId int64, title, imageUrl string) error {
	service, _ := oss.GetOssService()
	// 上传图片

	//视频文件上传
	err := service.UploadFileWithBytestream(conf.MinioConfig.VideoBucketName, reader, filename+".mp4", dataLen, videoContentType)
	if err != nil {
		log.Fatalln("OSS视频文件上传失败")
		return err
	}
	//抽取第一帧图片,得先上传视频文件，然后获取到视频文件的播放链接
	buffer, err := GetSnapshotImageBuffer(videoUrl, 1)
	if err != nil {
		log.Fatalln("视频封面图抽取失败！", err)
		return err
	}
	imgReader := strings.NewReader(buffer.String())
	_ = service.UploadFileWithBytestream(conf.MinioConfig.VideoBucketName, imgReader, filename+"-img.jpeg", int64(buffer.Len()), imageContentType)

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
	err = models.InsertVideo(&video)
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
	videoJson, _ := json.Marshal(video)
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
func GetVideoFeed(latestTime int64, nums int) ([]*model.Video, error, int64) {
	//获取latestTime时间之前的 不包括last？
	list, err := models.GetVideoFeedList(latestTime, nums)
	if err != nil {
		return nil, err, 0
	}
	resVideoList := make([]*model.Video, 0)
	resVideoList, err = convert.VideoSliceBo2Dto(list)
	if err != nil {
		return nil, err, 0
	}

	userVideoMap := map[int64]*model.User{}
	for i, item := range list {
		//如果这条视频已经查询过user了
		if v, ok := userVideoMap[item.ID]; ok {
			resVideoList[i].Author = v
			continue
		}
		//这一个还没被查询过
		userBo, _ := models.GetUserById(item.AuthorID)
		user, _ := convert.UserBo2Dto(*userBo)
		//缓存
		userVideoMap[item.ID] = user
		resVideoList[i].Author = user
	}
	return resVideoList, nil, list[len(list)-1].CreatedAt.UnixMilli()
}
