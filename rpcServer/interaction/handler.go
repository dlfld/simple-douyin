package main

import (
	"context"
	"github.com/douyin/common/crud"
	"github.com/douyin/kitex_gen/interaction"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
	"net/http"
	"time"
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
	resp = new(interaction.CommentActionResponse)
	actionType := req.ActionType
	if actionType == 1 { // 增加评论
		if req.CommentText == nil {
			resp.StatusCode = http.StatusOK
			msg := "增加评论时请输入评论内容"
			resp.StatusMsg = &msg
			return
		}
		m := models.Comment{
			VideoID: req.VideoId,
			UserID:  *req.UserId,
			Content: *req.CommentText,
		}
		_, err = Dao.InsertComment(&m)
		if err != nil {
			// TODO log err
			return
		}
		user, _ := Dao.SearchUserById(*req.UserId)
		resp.Comment = &model.Comment{
			Id:         0,
			User:       user,
			Content:    *req.CommentText,
			CreateDate: time.Now().String(),
		}
	} else if actionType == 2 { //删除评论
		if req.CommentId == nil {
			resp.StatusCode = http.StatusOK
			msg := "删除评论时请输入评论ID"
			resp.StatusMsg = &msg
			return
		}
		m := models.Comment{
			VideoID: req.VideoId,
			ID:      *req.CommentId,
		}
		_, err = Dao.DeleteComment(&m)
	} else {
		resp.StatusCode = http.StatusOK
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

// CommentList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) CommentList(ctx context.Context, req *interaction.CommentListRequest) (resp *interaction.CommentListResponse, err error) {
	dbList, err := Dao.SearchCommentListSort(req.VideoId)
	if err != nil {
		return nil, err
	}
	commentList := make([]*model.Comment, len(dbList))
	for i := 0; i < len(dbList); i++ {
		user, _ := Dao.SearchUserById(dbList[i].UserID)
		commentList[i] = &model.Comment{
			Id:         dbList[i].ID,
			User:       user,
			Content:    dbList[i].Content,
			CreateDate: dbList[i].CreateTime.String(),
		}
	}
	resp = &interaction.CommentListResponse{
		StatusCode:  http.StatusOK,
		CommentList: commentList,
	}
	return
}
