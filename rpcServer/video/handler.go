package main

import (
	"context"
	"github.com/douyin/kitex_gen/video"
	"github.com/douyin/models"
	"github.com/douyin/rpcServer/video/convert"
	"net/http"
	"strconv"
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
	return
}

// PublishList implements the VideoServiceImpl interface.
// 获取登录用户的视频发布列表，直接列出用户所有投稿过的视频。
func (s *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	// 根据登陆用户的id，查询用户所投稿过的所有视频
	videoList, err := models.FindVideoListBy("id", strconv.Itoa(int(req.GetUserId())))
	if err != nil {
		return nil, err
	}
	// 将从数据库中查询出来的bo对象转换为kitex的dto对象
	videoDtoSlice, err := convert.VideoSliceBo2Dto(videoList)
	if err != nil {
		return nil, err
	}
	// 封装返回结果
	resp = &video.PublishListResponse{VideoList: videoDtoSlice, StatusCode: http.StatusOK}
	return resp, err
}
