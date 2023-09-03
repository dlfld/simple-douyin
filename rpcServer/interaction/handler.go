package main

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/douyin/common/crud"
	"github.com/douyin/common/kafkaLog/productor"
	"github.com/douyin/kitex_gen/interaction"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
	d "github.com/douyin/rpcServer/interaction/dao"
)

var dao *d.Dao
var logger *productor.LogCollector

// InteractionServiceImpl implements the last service interface defined in the IDL.
type InteractionServiceImpl struct {
}

func InitDao() {
	dao = d.NewDao()
	dao.BloomFilter.PreLoadAll(dao.Mysql.GetCli())
}

func (s *InteractionServiceImpl) FavoriteAction(ctx context.Context, req *interaction.FavoriteActionRequest) (resp *interaction.FavoriteActionResponse, err error) {
	//if !dao.BloomFilter.IfVideoIdExists(req.VideoId) {
	//	return newFavoriteActionResp(-400, "入参无效"), errors.New("入参无效")
	//}
	m := models.FavoriteVideoRelation{
		VideoID: req.VideoId,
		UserID:  req.UserId,
	}
	authorId, err := dao.Mysql.SearchAuthorIdsByVideoId(req.VideoId)
	if err != nil {
		//logger.Error(fmt.Sprintf("FavoriteAction 执行错误[%v]", err))
		return newFavoriteActionResp(-500, "FavoriteAction 失败"), err
	}
	actionType := req.ActionType
	if actionType == 1 { // 点赞
		exists, _ := dao.Mysql.SearchFavoriteExist(&m)
		if exists {
			return newFavoriteActionResp(-400, "操作失败: 不能重复点赞"), errors.New("入参无效")
		}
		err = dao.Mysql.GetCli().Transaction(func(tx *gorm.DB) (err error) {
			_, err = dao.Mysql.InsertFavorite(&m)
			err = crud.FavoriteVideo(req.UserId, req.VideoId)
			_, err = dao.Mysql.VideoFavoriteCountIncr(req.VideoId, 1)
			_, err = dao.Mysql.UserFavoriteCountIncr(req.UserId, 1)
			_, err = dao.Mysql.UserTotalFavoritedCountIncr(authorId, 1)
			return err
		})
	} else if actionType == 2 { //取消点赞
		exists, _ := dao.Mysql.SearchFavoriteExist(&m)
		if !exists {
			//logger.Error(fmt.Sprintf("FavoriteAction 执行错误[%v]", err))
			return newFavoriteActionResp(-400, "操作失败: 您之前未点过赞, 无法取消点赞"), nil
		}
		err = dao.Mysql.GetCli().Transaction(func(tx *gorm.DB) (err error) {
			_, err = dao.Mysql.CancelFavorite(&m)
			err = crud.UnFavoriteVideo(req.UserId, req.VideoId)
			_, err = dao.Mysql.VideoFavoriteCountIncr(req.VideoId, -1)
			_, err = dao.Mysql.UserFavoriteCountIncr(req.UserId, -1)
			_, err = dao.Mysql.UserTotalFavoritedCountIncr(authorId, -1)
			return err
		})
	} else {
		return newFavoriteActionResp(-400, "actionType 输入错误：1-点赞，2-取消点赞"), nil
	}
	if err != nil {
		//logger.Error(fmt.Sprintf("FavoriteAction 执行错误[%v]", err))
		return newFavoriteActionResp(-500, "FavoriteAction 失败"), err
	}
	_ = dao.Redis.DelFavoriteVideoListByUserId(req.UserId)
	return newFavoriteActionResp(0, "操作成功"), nil
}

func (s *InteractionServiceImpl) FavoriteList(ctx context.Context, req *interaction.FavoriteListRequest) (resp *interaction.FavoriteListResponse, err error) {
	//if !dao.BloomFilter.IfUserIdExists(req.UserId) {
	//	return newFavoriteListResp(-400, "入参无效", nil), errors.New("入参无效")
	//}
	if videoList, err := dao.Redis.GetFavoriteVideoListByUserId(req.UserId); err == nil {
		return newFavoriteListResp(0, "ok", videoList), nil
	}
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
	_ = dao.Redis.SaveFavoriteVideoListByUserId(req.UserId, videoList)
	return newFavoriteListResp(0, "ok", videoList), err
}

// CommentAction implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) CommentAction(ctx context.Context, req *interaction.CommentActionRequest) (resp *interaction.CommentActionResponse, err error) {
	actionType := req.ActionType
	if actionType == 1 { // 增加评论
		//if !dao.BloomFilter.IfUserIdExists(*req.UserId) || !dao.BloomFilter.IfVideoIdExists(req.VideoId) {
		//	return newCommentActionResponse(-400, "入参无效", nil), errors.New("入参无效")
		//}
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
			//dao.BloomFilter.AddCommentIds(commentId)
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
		_ = dao.Redis.DelCommentListByVideoId(req.VideoId)
		return newCommentActionResponse(0, "增加评论成功", comment), nil
	} else if actionType == 2 { //删除评论
		//if req.CommentId == nil || !dao.BloomFilter.IfCommentIdExists(*req.CommentId) {
		//	return newCommentActionResponse(-500, "未输入commentId或无效", nil), errors.New("入参无效")
		//}
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
		_ = dao.Redis.DelCommentListByVideoId(req.VideoId)
		return newCommentActionResponse(0, "评论删除成功", nil), nil
	}
	return newCommentActionResponse(-400, "actionType 输入错误: 1-发布评论，2-删除评论", nil), nil
}

// CommentList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) CommentList(ctx context.Context, req *interaction.CommentListRequest) (resp *interaction.CommentListResponse, err error) {
	//if !dao.BloomFilter.IfVideoIdExists(req.VideoId) {
	//	return newCommentListResponse(-400, "入参无效", nil), errors.New("入参无效")
	//}
	if commentList, err := dao.Redis.GetCommentListByVideoId(req.VideoId); err == nil {
		return newCommentListResponse(0, "ok", commentList), nil
	}
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
	_ = dao.Redis.SaveCommentListByVideoId(req.VideoId, commentList)
	return newCommentListResponse(0, "ok", commentList), nil
}
