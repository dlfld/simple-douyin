package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/douyin/common/constant"
	"github.com/douyin/common/crud"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/kitex_gen/video"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	var feed []*model.Video
	var nextTime int64
	resp = new(video.FeedResponse)
	feed, nextTime, err = crud.GetUserFeed(req.GetUserId(), req.GetLatestTime())
	if err != nil {
		log.Fatalln("视频流调用失败")
		//往Kafka中写入错误日志
		LogCollector.Error(fmt.Sprintf("user[%d]:Failed to get video feed in %s, err=%s", req.GetUserId(), time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		// 返回给客户端错误信息
		constant.HandlerErr(constant.ErrFeedErr, resp)
		return resp, nil
	}
	statusMsg := "Success"
	// log.Println("%+v\n", feed)
	return &video.FeedResponse{VideoList: feed, StatusMsg: &statusMsg, StatusCode: 0, NextTime: &nextTime}, nil
}

// PublishAction implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishAction(ctx context.Context, req *video.PublishActionRequest) (resp *video.PublishActionResponse, err error) {
	reader := bytes.NewReader(req.GetData())
	fmt.Printf("data -> %+v\n", reader.Len())
	resp = new(video.PublishActionResponse)
	// 上传文件的文件名
	//filename := fmt.Sprintf("user-%d-%d", req.UserId, time.Now().Unix())
	userId := req.UserId
	log.Printf("userId --> %v\n", userId)
	title := req.Title
	dataLen := int64(len(req.GetData()))
	//执行视频上传逻辑
	if dataLen > 50*1000*1000 {
		constant.HandlerErr(constant.ErrVideoSizeMaxLimit, resp)
		return &video.PublishActionResponse{StatusCode: 1, StatusMsg: nil}, nil
	}
	if len(title) == 0 || len(title) > 50 {
		constant.HandlerErr(constant.ErrVideoTitleLength, resp)
		return &video.PublishActionResponse{StatusCode: 1, StatusMsg: nil}, nil
	}
	go func() {
		for i := 0; i < 10; i++ {
			err = UploadVideo(reader, dataLen, userId, title)
			if err != nil {
				//往Kafka中写入错误日志
				LogCollector.Error(fmt.Sprintf("user[%d]:Failed upload video in %s, err=%s", req.GetUserId(), time.Now().Format("2006-01-02 15:04:05"), err.Error()))
			}
		}

	}()
	statusMsg := "视频上传成功，后台上传完成之后便可查看"
	resp = &video.PublishActionResponse{StatusCode: 0, StatusMsg: &statusMsg}
	return resp, nil
}

// PublishList implements the VideoServiceImpl interface.
// 获取登录用户的视频发布列表，直接列出用户所有投稿过的视频。
func (s *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	// 根据登陆用户的id，查询用户所投稿过的所有视频
	resp = new(video.PublishListResponse)
	videoList, err := FindVideoListByUserId(int(req.GetUserId()))
	if err != nil {
		//往Kafka中写入错误日志
		LogCollector.Error(fmt.Sprintf("user[%d]:Failed to get user publish list in %s, err=%s", req.GetUserId(), time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		constant.HandlerErr(constant.ErrPublishList, resp)
		return resp, nil
	}
	// 封装返回结果
	resp = &video.PublishListResponse{VideoList: videoList, StatusCode: 0}
	return resp, err
}
