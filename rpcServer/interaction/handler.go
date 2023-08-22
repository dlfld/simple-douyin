package main

import (
	"context"
	"github.com/douyin/kitex_gen/interaction"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
	d "github.com/douyin/rpcServer/interaction/dao"
	"log"
	"net/http"
	"time"
)

var dao *d.Dao

// InteractionServiceImpl implements the last service interface defined in the IDL.
type InteractionServiceImpl struct {
	dao *d.Dao
}

func InitDao() {
	dao = d.NewDao()
}

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

func (s *InteractionServiceImpl) FavoriteAction(ctx context.Context, req *interaction.FavoriteActionRequest) (resp *interaction.FavoriteActionResponse, err error) {
	resp = new(interaction.FavoriteActionResponse)
	actionType := req.ActionType

	m := models.FavoriteVideoRelation{
		VideoID: req.VideoId,
		UserID:  req.UserId,
	}
	if actionType == 1 { // 点赞
		exists, _ := dao.Mysql.SearchFavoriteExist(&m)
		if exists {
			return newFavoriteActionResp(-400, "操作失败: 不能重复点赞"), nil
		}
		_, err = dao.Mysql.InsertFavorite(&m)
	} else if actionType == 2 { //取消点赞
		exists, _ := dao.Mysql.SearchFavoriteExist(&m)
		if !exists {
			return newFavoriteActionResp(-400, "操作失败: 您之前未点过赞, 无法取消点赞"), nil
		}
		_, err = dao.Mysql.CancelFavorite(&m)
	} else {
		return newFavoriteActionResp(-400, "actionType 错误"), nil
	}
	if err != nil {
		log.Println("FavoriteAction 执行错误")
		return newFavoriteActionResp(-500, "FavoriteAction 失败"), err
	}
	return newFavoriteActionResp(0, "操作成功"), nil
}

func (s *InteractionServiceImpl) FavoriteList(ctx context.Context, req *interaction.FavoriteListRequest) (resp *interaction.FavoriteListResponse, err error) {
	dbList, err := dao.Mysql.SearchVideoListById(req.UserId)
	if err != nil {
		return newFavoriteListResp(-500, "FavoriteList 错误", nil), err
	}
	authorIds := make([]int64, len(dbList))
	for i := 0; i < len(dbList); i++ {
		authorIds[i] = dbList[i].AuthorID
	}
	authorList, err := dao.Mysql.SearchUserByIds(authorIds, req.UserId)
	if err != nil || len(authorList) != len(dbList) {
		return newFavoriteListResp(-500, "FavoriteList 错误", nil), err
	}

	videoList := make([]*model.Video, len(dbList))
	for i := 0; i < len(dbList); i++ {
		videoList[i] = &model.Video{
			Id:            dbList[i].ID,
			Author:        authorList[i],
			PlayUrl:       dbList[i].PlayUrl,
			CoverUrl:      dbList[i].CoverUrl,
			FavoriteCount: dbList[i].FavoriteCount,
			CommentCount:  dbList[i].CommentCount,
			IsFavorite:    true,
			Title:         dbList[i].Title,
		}
	}

	return newFavoriteListResp(0, "ok", videoList), err
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
		_, err = dao.Mysql.InsertComment(&m)
		if err != nil {
			// TODO log err
			return
		}
		user, _ := dao.Mysql.SearchUserById(*req.UserId)
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
		_, err = dao.Mysql.DeleteComment(&m)
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
	dbList, err := dao.Mysql.SearchCommentListSort(req.VideoId)
	if err != nil {
		return nil, err
	}
	commentList := make([]*model.Comment, len(dbList))
	for i := 0; i < len(dbList); i++ {
		user, _ := dao.Mysql.SearchUserById(dbList[i].UserID)
		commentList[i] = &model.Comment{
			Id:         dbList[i].ID,
			User:       user,
			Content:    dbList[i].Content,
			CreateDate: dbList[i].CreatedTime.String(),
		}
	}
	resp = &interaction.CommentListResponse{
		StatusCode:  http.StatusOK,
		CommentList: commentList,
	}
	return
}
