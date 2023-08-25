package main

import (
	"context"
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/douyin/common/crud"
	"github.com/douyin/kitex_gen/interaction"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
	d "github.com/douyin/rpcServer/interaction/dao"
)

var dao *d.Dao

// InteractionServiceImpl implements the last service interface defined in the IDL.
type InteractionServiceImpl struct {
	dao *d.Dao
}

func InitDao() {
	dao = d.NewDao()
}

func (s *InteractionServiceImpl) FavoriteAction(ctx context.Context, req *interaction.FavoriteActionRequest) (resp *interaction.FavoriteActionResponse, err error) {
	m := models.FavoriteVideoRelation{
		VideoID: req.VideoId,
		UserID:  req.UserId,
	}
	authorId, err := dao.Mysql.SearchAuthorIdsByVideoId(req.VideoId)
	if err != nil {
		log.Println("FavoriteAction 执行错误")
		return newFavoriteActionResp(-500, "FavoriteAction 失败"), err
	}
	actionType := req.ActionType
	if actionType == 1 { // 点赞
		exists, _ := dao.Mysql.SearchFavoriteExist(&m)
		if exists {
			return newFavoriteActionResp(-400, "操作失败: 不能重复点赞"), nil
		}
		err = dao.Mysql.GetCli().Transaction(func(tx *gorm.DB) (err error) {
			_, err = dao.Mysql.InsertFavorite(&m)
			_, err = dao.Mysql.VideoFavoriteCountIncr(req.VideoId, 1)
			_, err = dao.Mysql.UserFavoriteCountIncr(req.UserId, 1)
			_, err = dao.Mysql.UserTotalFavoritedCountIncr(authorId, 1)
			return err
		})
	} else if actionType == 2 { //取消点赞
		exists, _ := dao.Mysql.SearchFavoriteExist(&m)
		if !exists {
			return newFavoriteActionResp(-400, "操作失败: 您之前未点过赞, 无法取消点赞"), nil
		}
		err = dao.Mysql.GetCli().Transaction(func(tx *gorm.DB) (err error) {
			_, err = dao.Mysql.CancelFavorite(&m)
			_, err = dao.Mysql.VideoFavoriteCountIncr(req.VideoId, -1)
			_, err = dao.Mysql.UserFavoriteCountIncr(req.UserId, -1)
			_, err = dao.Mysql.UserTotalFavoritedCountIncr(authorId, -1)
			return err
		})
	} else {
		return newFavoriteActionResp(-400, "actionType 输入错误：1-点赞，2-取消点赞"), nil
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

	// authorList, err := dao.Mysql.SearchUserByMids(authorIds, req.UserId)
	// if err != nil {
	// 	return newFavoriteListResp(-500, "FavoriteList 错误", nil), err
	// }
	// authorMap := make(map[int64]*model.User)
	// for _, v := range authorList {
	// 	authorMap[v.Id] = v
	// }
	authorMap, err := crud.GetAuthors(req.UserId, authorIds)

	videoList := make([]*model.Video, len(dbList))
	for i := 0; i < len(dbList); i++ {
		videoList[i] = &model.Video{
			Id:            dbList[i].ID,
			Author:        authorMap[authorIds[i]],
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
	actionType := req.ActionType
	if actionType == 1 { // 增加评论
		if req.CommentText == nil {
			return newCommentActionResponse(-500, "请输入评论内容", nil), err
		}
		m := models.Comment{
			VideoID:     req.VideoId,
			UserID:      *req.UserId,
			Content:     *req.CommentText,
			CreatedTime: time.Now(),
		}
		var commentId int64
		err = dao.Mysql.GetCli().Transaction(func(tx *gorm.DB) (err error) {
			commentId, err = dao.Mysql.InsertComment(&m)
			_, err = dao.Mysql.VideoCommentCountIncr(req.VideoId, 1)
			return err
		})
		user, _ := dao.Mysql.SearchUserByUserId(*req.UserId)
		if err != nil {
			return newCommentActionResponse(-500, "CommentAction 失败", nil), err
		}
		comment := &model.Comment{
			Id:         commentId,
			User:       user,
			Content:    *req.CommentText,
			CreateDate: m.CreatedTime.Format("01-02"),
		}
		return newCommentActionResponse(0, "ok", comment), nil
	} else if actionType == 2 { //删除评论
		if req.CommentId == nil {
			return newCommentActionResponse(-500, "删除评论时请输入评论ID", nil), err
		}
		m := models.Comment{
			ID: *req.CommentId,
		}
		err = dao.Mysql.GetCli().Transaction(func(tx *gorm.DB) (err error) {
			_, err = dao.Mysql.DeleteComment(&m)
			_, err = dao.Mysql.VideoCommentCountIncr(req.VideoId, -1)
			return err
		})
		if err != nil {
			return newCommentActionResponse(-500, "CommentAction 失败", nil), err
		}
		return newCommentActionResponse(0, "评论删除成功", nil), nil
	} else {
		return newCommentActionResponse(-400, "actionType 输入错误: 1-发布评论，2-删除评论", nil), nil
	}
}

// CommentList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) CommentList(ctx context.Context, req *interaction.CommentListRequest) (resp *interaction.CommentListResponse, err error) {
	dbList, err := dao.Mysql.SearchCommentListSort(req.VideoId)
	if err != nil {
		return newCommentListResponse(-500, "CommentList 失败", nil), err
	}
	commentUserIds := make([]int64, len(dbList))
	for i := 0; i < len(dbList); i++ {
		commentUserIds[i] = dbList[i].UserID
	}
	commentUserList, err := dao.Mysql.SearchUserByMids(commentUserIds, *req.UserId)
	if err != nil {
		return newCommentListResponse(-500, "CommentList 失败", nil), err
	}
	commentUserMap := make(map[int64]*model.User)
	for _, v := range commentUserList {
		commentUserMap[v.Id] = v
	}
	commentList := make([]*model.Comment, len(dbList))
	for i := 0; i < len(dbList); i++ {
		commentList[i] = &model.Comment{
			Id:         dbList[i].ID,
			User:       commentUserMap[commentUserIds[i]],
			Content:    dbList[i].Content,
			CreateDate: dbList[i].CreatedTime.Format("01-02"),
		}
	}
	return newCommentListResponse(0, "ok", commentList), nil
}
