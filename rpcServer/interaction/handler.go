package main

import (
	"context"
	"github.com/douyin/kitex_gen/model"
	"net/http"

	"github.com/douyin/common/crud"
	"github.com/douyin/kitex_gen/interaction"
	"github.com/douyin/models"
)

var Dao *crud.CachedCRUD

// InteractionServiceImpl implements the last service interface defined in the IDL.
type InteractionServiceImpl struct{}

func InitDao() (err error) {
	Dao, err = crud.NewCachedCRUD()
	return
}

// FavoriteAction implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) FavoriteAction(ctx context.Context, req *interaction.FavoriteActionRequest) (resp *interaction.FavoriteActionResponse, err error) {
	resp = new(interaction.FavoriteActionResponse)
	actionType := req.ActionType

	m := models.FavoriteVideoRelation{
		VideoID: req.VideoId,
		UserID:  req.UserId,
	}
	if actionType == 1 { // 点赞
		exist, _ := Dao.SearchFavoriteExist(&m)
		if exist {
			resp.StatusCode = http.StatusOK
			msg := "该记录已经存在"
			resp.StatusMsg = &msg
			return
		}
		_, err = Dao.InsertFavorite(&m)
	} else if actionType == 2 { //取消点赞
		_, err = Dao.CancelFavorite(&m)
	} else {
		resp.StatusCode = http.StatusInternalServerError
		msg := "actionType 错误"
		resp.StatusMsg = &msg
		return
	}
	if err != nil {
		// TODO: log err
	}
	resp.StatusCode = http.StatusOK
	msg := "ok"
	resp.StatusMsg = &msg
	return
}

// FavoriteList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) FavoriteList(ctx context.Context, req *interaction.FavoriteListRequest) (resp *interaction.FavoriteListResponse, err error) {
	resp = new(interaction.FavoriteListResponse)
	dbList, err := Dao.SearchVideoListById(req.UserId)
	if err != nil {
		return nil, err
	}
	videoList := make([]*model.Video, len(dbList))
	for i := 0; i < len(dbList); i++ {
		author, _ := Dao.SearchUserById(dbList[i].AuthorID)
		m := models.FavoriteVideoRelation{
			VideoID: dbList[i].ID,
			UserID:  req.UserId,
		}
		isFavorite, _ := Dao.SearchFavoriteExist(&m)
		videoList[i] = &model.Video{
			Id:            dbList[i].ID,
			Author:        author,
			PlayUrl:       dbList[i].PlayUrl,
			CoverUrl:      dbList[i].CoverUrl,
			FavoriteCount: dbList[i].FavoriteCount,
			CommentCount:  dbList[i].CommentCount,
			IsFavorite:    isFavorite,
			Title:         dbList[i].Title,
		}
	}
	resp.StatusCode = http.StatusOK
	msg := "ok"
	resp.StatusMsg = &msg
	resp.VideoList = videoList
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
