package main

import (
	"context"
	"fmt"
	"github.com/douyin/common/bloom"
	"github.com/douyin/common/constant"

	"github.com/douyin/common/crud"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/kitex_gen/relation"
	"github.com/douyin/models"
)

var bf *bloom.Filter

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

func usersToKitex(self uint, users []*models.User) (kitexList []*model.User) {
	for _, v := range users {
		Id := v.ID
		Name := v.UserName
		FollowCount := int64(v.FollowingCount)
		FollowerCount := int64(v.FollowerCount)
		IsFollow := crud.IsFollow(self, v.ID)
		Avatar := v.Avatar
		BackgroundImage := v.BackgroundImage
		Signature := v.Signature
		TotalFavorited := int64(v.TotalFavorited)
		WorkCount := int64(v.WorkCount)
		FavoriteCount := int64(v.FavoriteCount)
		u := model.User{
			Id:              int64(Id),
			Name:            Name,
			FollowCount:     &FollowCount,
			FollowerCount:   &FollowerCount,
			IsFollow:        IsFollow,
			Avatar:          &Avatar,
			BackgroundImage: &BackgroundImage,
			Signature:       &Signature,
			TotalFavorited:  &TotalFavorited,
			WorkCount:       &WorkCount,
			FavoriteCount:   &FavoriteCount,
		}
		kitexList = append(kitexList, &u)
	}
	return
}

func friendUsersToKitex(friendUsers []*models.User) (kitexList []*model.FriendUser) {
	for _, v := range friendUsers {
		Id := v.ID
		Name := v.UserName
		FollowCount := int64(v.FollowingCount)
		FollowerCount := int64(v.FollowerCount)
		IsFollow := true
		Avatar := v.Avatar
		BackgroundImage := v.BackgroundImage
		Signature := v.Signature
		TotalFavorited := int64(v.TotalFavorited)
		WorkCount := int64(v.WorkCount)
		FavoriteCount := int64(v.FavoriteCount)
		u := model.FriendUser{
			Id:              int64(Id),
			Name:            Name,
			FollowCount:     &FollowCount,
			FollowerCount:   &FollowerCount,
			IsFollow:        IsFollow,
			Avatar:          &Avatar,
			BackgroundImage: &BackgroundImage,
			Signature:       &Signature,
			TotalFavorited:  &TotalFavorited,
			WorkCount:       &WorkCount,
			FavoriteCount:   &FavoriteCount,
		}
		kitexList = append(kitexList, &u)
	}
	return
}

// FollowAction implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowAction(ctx context.Context, req *relation.FollowActionRequest) (resp *relation.FollowActionResponse, err error) {
	// TODO: Your code here...
	var msg string
	// crud, _ := crud.NewCachedCRUD()
	resp = new(relation.FollowActionResponse)
	resp.StatusMsg = &msg
	userId := uint(req.FromUserId)

	exists, err := bf.CheckIfUserIdExists(req.ToUserId)
	if err != nil {
		logCollector.Error(fmt.Sprintf("User bloom_user err[%v]", err))
	} else {
		if !exists {
			constant.HandlerErr(constant.ErrBloomUser, resp)
			return resp, nil
		}
	}

	switch req.ActionType {
	case 1:
		// err = models.Follow(userId, uint(req.ToUserId))
		err = crud.RelationFollow(userId, uint(req.ToUserId))
		if err != nil {
			// msg = err.Error()
			logCollector.Error(fmt.Sprintf("follow action error: %v", err))
		} else {
			crud.DeletePublishListCache(int(userId))
			msg = "follow ok"
		}

	case 2:
		// err = models.Unfollow(userId, uint(req.ToUserId))
		err = crud.RelationUnFollow(userId, uint(req.ToUserId))
		if err != nil {
			logCollector.Error(fmt.Sprintf("unfollow action error: %v", err))
		} else {
			crud.DeletePublishListCache(int(userId))
			msg = "unfollow ok"
		}

	default:
		msg = fmt.Sprintf("unknow action type: %d", req.ActionType)
	}
	return resp, nil
}

// FollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowList(ctx context.Context, req *relation.FollowingListRequest) (resp *relation.FollowingListResponse, err error) {
	// TODO: Your code here...
	resp = new(relation.FollowingListResponse)
	exists, err := bf.CheckIfUserIdExists(req.UserId)
	if err != nil {
		logCollector.Error(fmt.Sprintf("User bloom_user err[%v]", err))
	} else {
		if !exists {
			constant.HandlerErr(constant.ErrBloomUser, resp)
			return resp, nil
		}
	}
	var msg string = "get follow list ok"
	// crud, _ := crud.NewCachedCRUD()
	// userList, _ := models.GetFollowList(uint(req.UserId))
	userList, err := crud.RelationGetFollows(uint(req.UserId))
	if err != nil {
		logCollector.Error(fmt.Sprintf("get follow list error: %v", err))
	}
	kitexList := usersToKitex(uint(req.UserId), userList)

	return &relation.FollowingListResponse{StatusCode: 0, StatusMsg: &msg, UserList: kitexList}, nil
}

// FollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowerList(ctx context.Context, req *relation.FollowerListRequest) (resp *relation.FollowerListResponse, err error) {
	// TODO: Your code here...
	resp = new(relation.FollowerListResponse)
	exists, err := bf.CheckIfUserIdExists(req.UserId)
	if err != nil {
		logCollector.Error(fmt.Sprintf("User bloom_user err[%v]", err))
	} else {
		if !exists {
			constant.HandlerErr(constant.ErrBloomUser, resp)
			return resp, nil
		}
	}

	var msg string = "get follow list ok"
	// crud, _ := crud.NewCachedCRUD()
	// var kitexList []*model.User
	// userList, _ := models.GetFollowerList(uint(req.UserId))
	userList, err := crud.RelationGetFollowers(uint(req.UserId))
	if err != nil {
		logCollector.Error(fmt.Sprintf("get follower list error: %v", err))
	}
	kitexList := usersToKitex(uint(req.UserId), userList)

	return &relation.FollowerListResponse{StatusCode: 0, StatusMsg: &msg, UserList: kitexList}, nil
}

// FriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FriendList(ctx context.Context, req *relation.RelationFriendListRequest) (resp *relation.RelationFriendListResponse, err error) {
	// TODO: Your code here...
	resp = new(relation.RelationFriendListResponse)
	exists, err := bf.CheckIfUserIdExists(req.UserId)
	if err != nil {
		logCollector.Error(fmt.Sprintf("User bloom_user err[%v]", err))
	} else {
		if !exists {
			constant.HandlerErr(constant.ErrBloomUser, resp)
			return resp, nil
		}
	}

	var msg string = "get follow list ok"
	// crud, _ := crud.NewCachedCRUD()
	userList, err := crud.RelationGetFriends(uint(req.UserId))
	if err != nil {
		logCollector.Error(fmt.Sprintf("get friend list error: %v", err))
	}
	// userList, _ := models.GetFriendList(uint(req.UserId))
	// var kitexList []*model.FriendUser
	kitexList := friendUsersToKitex(userList)

	return &relation.RelationFriendListResponse{StatusCode: 0, StatusMsg: &msg, UserList: kitexList}, nil
}
