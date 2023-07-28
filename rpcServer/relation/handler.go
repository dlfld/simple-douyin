package main

import (
	"context"
	"github.com/douyin/kitex_gen/relation"
)

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

// FollowAction implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowAction(ctx context.Context, req *relation.FollowActionRequest) (resp *relation.FollowActionResponse, err error) {
	// TODO: Your code here...
	return
}

// FollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowList(ctx context.Context, req *relation.FollowingListRequest) (resp *relation.FollowingListResponse, err error) {
	// TODO: Your code here...
	return
}

// FollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowerList(ctx context.Context, req *relation.FollowerListRequest) (resp *relation.FollowerListResponse, err error) {
	// TODO: Your code here...
	return
}

// FriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FriendList(ctx context.Context, req *relation.RelationFriendListRequest) (resp *relation.RelationFriendListResponse, err error) {
	// TODO: Your code here...
	return
}
