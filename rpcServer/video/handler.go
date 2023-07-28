package main

import (
	"context"
	"github.com/douyin/kitex_gen/video"
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
func (s *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	// TODO: Your code here...
	return
}
