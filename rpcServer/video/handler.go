package main

import (
	"bytes"
	"context"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/oss"
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
	service, _ := oss.GetOssService()
	reader := bytes.NewReader(req.GetData())
	// 上传文件的文件名
	filename := req.GetTitle()
	// TODO 根据Token获取用户信息，然后根据用户信息写入用户投稿的视频，在redis中加入这一条视频
	//videoUrl := fmt.Sprintf("http://%s/%s/%s", conf.MinioConfig.IP, conf.MinioConfig.VideoBucketName, filename)

	// TODO 魔法值需要改
	contentType := "application/mp4"
	err = service.UploadFileWithBytestream(conf.MinioConfig.VideoBucketName, reader, filename, int64(len(req.GetData())), contentType)
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
