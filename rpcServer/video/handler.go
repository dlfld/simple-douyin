package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/douyin/common/conf"
	video2 "github.com/douyin/common/crud/video"
	"github.com/douyin/kitex_gen/video"
	"net/http"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	// TODO: Your code here...
	return
}

// PublishAction implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishAction(ctx context.Context, req *video.PublishActionRequest) (resp *video.PublishActionResponse, err error) {
	// TODO: Your code here...

	reader := bytes.NewReader(req.GetData())
	// 上传文件的文件名
	filename := req.GetTitle()
	userId := req.UserId
	// TODO 根据Token获取用户信息，然后根据用户信息写入用户投稿的视频，在redis中加入这一条视频
	videoUrl := fmt.Sprintf("http://%s/%s/%s", conf.MinioConfig.IP, conf.MinioConfig.VideoBucketName, filename)
	// TODO 魔法值需要改
	contentType := "application/mp4"
	dataLen := int64(len(req.GetData()))
	//执行视频上传逻辑
	err = video2.UploadVideo(reader, filename, contentType, videoUrl, dataLen, userId)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	// TODO 魔法值需要改
	statusMsg := "success"
	resp = &video.PublishActionResponse{StatusCode: http.StatusOK, StatusMsg: &statusMsg}
	return resp, nil
}

// PublishList implements the VideoServiceImpl interface.
// 获取登录用户的视频发布列表，直接列出用户所有投稿过的视频。
func (s *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	// 根据登陆用户的id，查询用户所投稿过的所有视频
	videoList, err := video2.FindVideoListByUserId(int(req.GetUserId()))
	if err != nil {
		return nil, err
	}
	// 封装返回结果
	resp = &video.PublishListResponse{VideoList: videoList, StatusCode: http.StatusOK}
	return resp, err
}
