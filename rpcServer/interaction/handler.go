package main

import (
	"context"
	"github.com/douyin/common/crud"
	"github.com/douyin/kitex_gen/interaction"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
)

var dao *crud.CachedCRUD

// InteractionServiceImpl implements the last service interface defined in the IDL.
type InteractionServiceImpl struct{}

// FavoriteAction implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) FavoriteAction(ctx context.Context, req *interaction.FavoriteActionRequest) (resp *interaction.FavoriteActionResponse, err error) {
	// TODO: Your code here...
	m := models.FavoriteVideoRelation{
		VideoID: req.VideoId,
		UserID:  req.,
	}

	return
}

// FavoriteList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) FavoriteList(ctx context.Context, req *interaction.FavoriteListRequest) (resp *interaction.FavoriteListResponse, err error) {
	resp = &interaction.FavoriteListResponse{
		StatusCode: 123,
		StatusMsg:  new(string),
		VideoList: []*model.Video{
			&model.Video{
				Id: 12,
				Author: &model.User{
					Id:   1234,
					Name: "naruto",
				},
			},
			&model.Video{
				Id: 44,
				Author: &model.User{
					Id:   567,
					Name: "kakasi",
				},
			},
		},
	}
	return
}

// CommentAction implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) CommentAction(ctx context.Context, req *interaction.CommentActionRequest) (resp *interaction.CommentActionResponse, err error) {
	// TODO: Your code here...
	return
}

// CommentList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) CommentList(ctx context.Context, req *interaction.CommentListRequest) (resp *interaction.CommentListResponse, err error) {
	// TODO: Your code here...
	return
}
