package main

import (
	"github.com/douyin/kitex_gen/interaction"
	"github.com/douyin/kitex_gen/model"
)

func newFavoriteActionResp(code int32, msg string) *interaction.FavoriteActionResponse {
	return &interaction.FavoriteActionResponse{
		StatusCode: code,
		StatusMsg:  &msg,
	}
}

func newFavoriteListResp(code int32, msg string, vlist []*model.Video) *interaction.FavoriteListResponse {
	return &interaction.FavoriteListResponse{
		StatusCode: code,
		StatusMsg:  &msg,
		VideoList:  vlist,
	}
}

func newCommentActionResponse(code int32, msg string, comment *model.Comment) *interaction.CommentActionResponse {
	return &interaction.CommentActionResponse{
		StatusCode: code,
		StatusMsg:  &msg,
		Comment:    comment,
	}
}
