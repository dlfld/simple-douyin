package main

import (
	"context"

	"github.com/douyin/common/crud"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/kitex_gen/relation"
	"github.com/douyin/models"
)

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

func usersToKitex(users []*models.User) (kitexList []*model.User) {
	for _, v := range users {
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
	crud, _ := crud.NewCachedCRUD()
	resp = new(relation.FollowActionResponse)
	resp.StatusMsg = &msg
	// var token2UserID = map[string]uint{
	// 	"token1": 1,
	// 	"token2": 2,
	// }
	// userId, has := token2UserID[req.Token]
	// if !has {
	// 	msg = "Token error"
	// 	return
	// }
	userId := uint(req.FromUserId)

	switch req.ActionType {
	case 1:
		// err = models.Follow(userId, uint(req.ToUserId))
		err = crud.RelationFollow(userId, uint(req.ToUserId))
		if err != nil {
			msg = err.Error()
		} else {
			msg = "follow ok"
		}
	case 2:
		// err = models.Unfollow(userId, uint(req.ToUserId))
		err = crud.RelationUnFollow(userId, uint(req.ToUserId))
		if err != nil {
			msg = err.Error()
		} else {
			msg = "unfollow ok"
		}

	default:
		msg = "unknow action type"
	}
	return
}

// FollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowList(ctx context.Context, req *relation.FollowingListRequest) (resp *relation.FollowingListResponse, err error) {
	// TODO: Your code here...
	var msg string = "get follow list ok"
	crud, _ := crud.NewCachedCRUD()
	// userList, _ := models.GetFollowList(uint(req.UserId))
	userList, _ := crud.RelationGetFollows(uint(req.UserId))
	kitexList := usersToKitex(userList)

	return &relation.FollowingListResponse{StatusCode: 0, StatusMsg: &msg, UserList: kitexList}, err
}

// FollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowerList(ctx context.Context, req *relation.FollowerListRequest) (resp *relation.FollowerListResponse, err error) {
	// TODO: Your code here...
	var msg string = "get follow list ok"
	crud, _ := crud.NewCachedCRUD()
	// var kitexList []*model.User
	// userList, _ := models.GetFollowerList(uint(req.UserId))
	userList, _ := crud.RelationGetFollowers(uint(req.UserId))
	kitexList := usersToKitex(userList)

	return &relation.FollowerListResponse{StatusCode: 0, StatusMsg: &msg, UserList: kitexList}, err
}

// FriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FriendList(ctx context.Context, req *relation.RelationFriendListRequest) (resp *relation.RelationFriendListResponse, err error) {
	// TODO: Your code here...
	var msg string = "get follow list ok"
	crud, _ := crud.NewCachedCRUD()
	userList, _ := crud.RelationGetFriends(uint(req.UserId))
	// userList, _ := models.GetFriendList(uint(req.UserId))
	// var kitexList []*model.FriendUser
	kitexList := friendUsersToKitex(userList)

	return &relation.RelationFriendListResponse{StatusCode: 0, StatusMsg: &msg, UserList: kitexList}, err
}
