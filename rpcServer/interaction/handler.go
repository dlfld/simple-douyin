package main

import (
	"context"
	"fmt"
	"github.com/douyin/common/constant"
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
var logCollector *productor.LogCollector

// InteractionServiceImpl implements the last service interface defined in the IDL.
type InteractionServiceImpl struct {
}

func InitDao() {
	dao = d.NewDao()
	dao.BloomFilter.PreLoadAll(dao.Mysql.GetCli())
}

func (s *InteractionServiceImpl) FavoriteAction(ctx context.Context, req *interaction.FavoriteActionRequest) (resp *interaction.FavoriteActionResponse, err error) {
	resp = new(interaction.FavoriteActionResponse)
	exists, err := dao.BloomFilter.CheckIfVideoIdExists(req.VideoId)
	if err != nil {
		logCollector.Error(fmt.Sprintf("Interaction bloom_video err[%v]", err))
	} else {
		if !exists {
			constant.HandlerErr(constant.ErrBloomVideo, resp)
			return resp, nil
		}
	}
	m := models.FavoriteVideoRelation{
		VideoID: req.VideoId,
		UserID:  req.UserId,
	}
	authorId, err := dao.Mysql.SearchAuthorIdsByVideoId(req.VideoId)
	if err != nil {
		//logger.Error(fmt.Sprintf("FavoriteAction 执行错误[%v]", err))
		constant.HandlerErr(constant.ErrFavoriteAction, resp)
		return resp, nil
	}
	actionType := req.ActionType
	if actionType == 1 { // 点赞
		exists, _ := dao.Mysql.SearchFavoriteExist(&m)
		if exists {
			constant.HandlerErr(constant.ErrFavoriteAction, resp)
			return resp, nil
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
			constant.HandlerErr(constant.ErrFavoriteAction, resp)
			return resp, nil
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
		constant.HandlerErr(constant.ErrFavoriteAction, resp)
		return resp, nil
	}
	if err != nil {
		//logger.Error(fmt.Sprintf("FavoriteAction 执行错误[%v]", err))
		constant.HandlerErr(constant.ErrFavoriteAction, resp)
		return resp, nil
	}
	_ = dao.Redis.DelFavoriteVideoListByUserId(req.UserId)
	return newFavoriteActionResp(0, "操作成功"), nil
}

func (s *InteractionServiceImpl) FavoriteList(ctx context.Context, req *interaction.FavoriteListRequest) (resp *interaction.FavoriteListResponse, err error) {
	resp = new(interaction.FavoriteListResponse)
	exists, err := dao.BloomFilter.CheckIfUserIdExists(req.UserId)
	if err != nil {
		logCollector.Error(fmt.Sprintf("Interaction bloom_user err[%v]", err))
	} else {
		if !exists {
			constant.HandlerErr(constant.ErrBloomUser, resp)
			return resp, nil
		}
	}

	if videoList, err := dao.Redis.GetFavoriteVideoListByUserId(req.UserId); err == nil {
		return newFavoriteListResp(0, "ok", videoList), nil
	}
	dbList, err := dao.Mysql.SearchVideoListById(req.UserId)
	if err != nil {
		constant.HandlerErr(constant.ErrFavoriteList, resp)
		return resp, nil
	}
	authorIds := make([]int64, len(dbList))
	for i := 0; i < len(dbList); i++ {
		authorIds[i] = dbList[i].AuthorID
	}

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
	resp = new(interaction.CommentActionResponse)
	exists, err := dao.BloomFilter.CheckIfVideoIdExists(req.VideoId)
	if err != nil {
		logCollector.Error(fmt.Sprintf("Interaction bloom_video err[%v]", err))
	} else {
		if !exists {
			constant.HandlerErr(constant.ErrBloomVideo, resp)
			return resp, nil
		}
	}
	actionType := req.ActionType
	if actionType != 1 && actionType != 2 {
		constant.HandlerErr(constant.ErrCommentAction, resp)
		return resp, nil
	}
	if actionType == 1 { // 增加评论
		if req.CommentText == nil {
			constant.HandlerErr(constant.ErrCommentAction, resp)
			return resp, nil
		}

		exists, err = dao.BloomFilter.CheckIfVideoIdExists(req.VideoId)
		if err != nil {
			logCollector.Error(fmt.Sprintf("Interaction bloom_video err[%v]", err))
		} else {
			if !exists {
				constant.HandlerErr(constant.ErrBloomVideo, resp)
				return resp, nil
			}
		}

		exists, err = dao.BloomFilter.CheckIfUserIdExists(*req.UserId)
		if err != nil {
			logCollector.Error(fmt.Sprintf("Interaction bloom_user err[%v]", err))
		} else {
			if !exists {
				constant.HandlerErr(constant.ErrBloomUser, resp)
				return resp, nil
			}
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
			constant.HandlerErr(constant.ErrCommentAction, resp)
			return resp, nil
		}
		comment := &model.Comment{
			Id:         commentId,
			User:       user,
			Content:    *req.CommentText,
			CreateDate: m.CreatedTime.Format("01-02"),
		}
		_ = dao.Redis.DelCommentListByVideoId(req.VideoId)
		dao.BloomFilter.AddCommentId(commentId)
		return newCommentActionResponse(0, "增加评论成功", comment), nil
	} else if actionType == 2 { //删除评论
		if req.CommentId == nil {
			constant.HandlerErr(constant.ErrCommentAction, resp)
			return resp, nil
		}

		exists, err = dao.BloomFilter.CheckIfCommentIdExists(*req.CommentId)
		if err != nil {
			logCollector.Error(fmt.Sprintf("Interaction bloom_comment err[%v]", err))
		} else {
			if !exists {
				constant.HandlerErr(constant.ErrBloomComment, resp)
				return resp, nil
			}
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
			constant.HandlerErr(constant.ErrCommentAction, resp)
			return resp, nil
		}
		_ = dao.Redis.DelCommentListByVideoId(req.VideoId)
		return newCommentActionResponse(0, "评论删除成功", nil), nil
	}
	return
}

// CommentList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) CommentList(ctx context.Context, req *interaction.CommentListRequest) (resp *interaction.CommentListResponse, err error) {
	resp = new(interaction.CommentListResponse)
	exists, err := dao.BloomFilter.CheckIfVideoIdExists(req.VideoId)
	if err != nil {
		logCollector.Error(fmt.Sprintf("Interaction bloom_video err[%v]", err))
	} else {
		if !exists {
			constant.HandlerErr(constant.ErrBloomVideo, resp)
			return resp, nil
		}
	}

	if commentList, err := dao.Redis.GetCommentListByVideoId(req.VideoId); err == nil {
		return newCommentListResponse(0, "ok", commentList), nil
	}
	dbList, err := dao.Mysql.SearchCommentListSort(req.VideoId)
	if err != nil {
		constant.HandlerErr(constant.ErrCommentList, resp)
		return resp, nil
	}
	commentUserIds := make([]int64, len(dbList))
	for i := 0; i < len(dbList); i++ {
		commentUserIds[i] = dbList[i].UserID
	}
	commentUserList, err := dao.Mysql.SearchUserByMids(commentUserIds, *req.UserId)
	if err != nil {
		constant.HandlerErr(constant.ErrCommentList, resp)
		return resp, nil
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
