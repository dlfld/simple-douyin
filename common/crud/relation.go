package crud

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
	"github.com/go-redis/redis/v8"
	"github.com/u2takey/go-utils/pointer"
)

func userRelationFollowKey(userID uint) string {
	return fmt.Sprintf("relation:follow:%d", userID)
}
func userRelationFollowerKey(userID uint) string {
	return fmt.Sprintf("relation:follower:%d", userID)
}

// RelationFollow 在缓存中建立关注关系
func RelationFollow(userID, toUserID uint) (err error) {
	if userID == toUserID {
		return
	}
	pipline := crud.redis.Pipeline()
	defer pipline.Close()
	models.Follow(userID, toUserID)
	// 将被关注用户添加到关注列表
	res := pipline.SAdd(context.Background(), userRelationFollowKey(userID), toUserID)
	if res.Err() != nil {
		return res.Err()
	}

	// 将关注用户添加到粉丝列表
	res = pipline.SAdd(context.Background(), userRelationFollowerKey(toUserID), userID)
	if res.Err() != nil {
		return res.Err()
	}
	// 执行缓存操作
	_, err = pipline.Exec(context.Background())
	return err
}

// RelationUnFollow 在缓存中取消关注关系
func RelationUnFollow(userID, toUserID uint) (err error) {
	if userID == toUserID {
		return
	}
	pipline := crud.redis.Pipeline()
	defer pipline.Close()
	models.Unfollow(userID, toUserID)
	// 从关注列表中移除被取消关注用户
	res := pipline.SRem(context.Background(), userRelationFollowKey(userID), toUserID)
	if res.Err() != nil {
		return res.Err()
	}

	// 从粉丝列表中移除取消关注用户
	res = pipline.SRem(context.Background(), userRelationFollowerKey(toUserID), userID)
	if res.Err() != nil {
		return res.Err()
	}
	// 执行缓存操作
	_, err = pipline.Exec(context.Background())
	return err
}

// RelationGetFollows 获取用户关注列表
func RelationGetFollows(userID uint) (userList []*models.User, err error) {
	ex := crud.redis.Exists(context.Background(), userRelationFollowKey(userID))
	if ex.Val() == 0 {
		var users []*models.User
		users, err = models.GetFollowList(userID)
		CacheRelationFollows(userID, users)
		CacheUsersInfo(users)
		return users, err
	}
	// 获取用户关注列表
	res := crud.redis.SMembers(context.Background(), userRelationFollowKey(userID))
	if res.Err() != nil {
		return nil, res.Err()
	}
	// 获取关注用户信息
	var ids = res.Val()
	userList, err = GetUsersByID(ids)
	return userList, err
}

func CacheRelationFollows(userID uint, follows []*models.User) {

	pipline := crud.redis.Pipeline()
	defer pipline.Close()
	for _, user := range follows {
		pipline.SAdd(context.Background(), userRelationFollowKey(userID), user.ID)
	}
	pipline.Expire(context.Background(), userRelationFollowKey(userID), time.Hour)
	pipline.Exec(context.Background())
}

func CacheRelationFollowers(userID uint, followers []*models.User) {
	pipline := crud.redis.Pipeline()
	defer pipline.Close()

	for _, user := range followers {
		pipline.SAdd(context.Background(), userRelationFollowerKey(userID), user.ID)
	}

	pipline.Expire(context.Background(), userRelationFollowerKey(userID), time.Hour)

	pipline.Exec(context.Background())
}

// RelationGetFollowers 获取用户的粉丝列表
func RelationGetFollowers(userID uint) (userList []*models.User, err error) {
	ex := crud.redis.Exists(context.Background(), userRelationFollowerKey(userID))
	if ex.Val() == 0 {
		var users []*models.User
		users, _ = models.GetFollowerList(userID)
		CacheRelationFollowers(userID, users)
		CacheUsersInfo(users)
	}
	ex = crud.redis.Exists(context.Background(), userRelationFollowKey(userID))
	if ex.Val() == 0 {
		var users []*models.User
		users, _ = models.GetFollowList(userID)
		CacheRelationFollows(userID, users)
		CacheUsersInfo(users)
	}
	// 获取用户的粉丝列表
	res := crud.redis.SMembers(context.Background(), userRelationFollowerKey(userID))
	if res.Err() != nil {
		return nil, res.Err()
	}
	var ids = res.Val()
	// 获取粉丝用户信息
	userList, err = GetUsersByID(ids)
	return userList, err
}

// RelationGetFriends 获取用户的好友列表
func RelationGetFriends(userID uint) (userList []*models.User, err error) {
	ex := crud.redis.Exists(context.Background(), userRelationFollowerKey(userID))
	if ex.Val() == 0 {
		var users []*models.User
		users, err = models.GetFollowerList(userID)
		CacheRelationFollowers(userID, users)
		CacheUsersInfo(users)
		return users, err
	}
	// 获取交集，即用户的好友列表
	res := crud.redis.SInter(context.Background(), userRelationFollowerKey(userID), userRelationFollowKey(userID))
	if res.Err() != nil {
		return nil, res.Err()
	}
	var ids = res.Val()
	// 获取好友用户信息
	userList, err = GetUsersByID(ids)
	return userList, err
}

// CacheUserInfo 将用户信息存入缓存
func CacheUserInfo(user *models.User) (err error) {
	// 序列化用户信息
	data, err := sonic.Marshal(user)
	if err != nil {
		return err
	}
	// 将用户信息存入缓存的哈希表中
	res := crud.redis.HSet(context.Background(), "UserInfoCache", fmt.Sprintf("%d", user.ID), string(data))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

// CacheUsersInfo 批量将用户信息存入缓存
func CacheUsersInfo(users []*models.User) (err error) {
	pipline := crud.redis.Pipeline()
	// 序列化用户信息
	for _, v := range users {
		data, err := sonic.Marshal(v)
		if err != nil {
			return err
		}
		res := pipline.HSet(
			context.Background(),
			"UserInfoCache",
			fmt.Sprintf("%d", v.ID),
			string(data),
		)
		if res.Err() != nil {
			return res.Err()
		}
	}
	// 将用户信息存入缓存的哈希表中

	_, err = pipline.Exec(context.Background())
	return
}

// GetAuthor 获取用户信息
func GetAuthor(self uint, UserID uint) (user *model.User, err error) {
	ormmodel, _ := GetUserInfo(strconv.Itoa(int(UserID)))
	user = &model.User{
		Id:              int64(ormmodel.ID),
		Name:            ormmodel.UserName,
		FollowCount:     pointer.Int64Ptr(int64(ormmodel.FollowingCount)),
		FollowerCount:   pointer.Int64Ptr(int64(ormmodel.FollowerCount)),
		Avatar:          pointer.StringPtr(ormmodel.Avatar),
		IsFollow:        IsFollow(self, UserID),
		BackgroundImage: pointer.StringPtr(ormmodel.BackgroundImage),
		FavoriteCount:   pointer.Int64Ptr(int64(ormmodel.FavoriteCount)),
		TotalFavorited:  pointer.Int64Ptr(int64(ormmodel.TotalFavorited)),
		WorkCount:       pointer.Int64Ptr(int64(ormmodel.WorkCount)),
	}
	return
}

func GetAuthors(self int64, UserIDs []int64) (users map[int64]*model.User, err error) {
	userIDStrs := make([]string, len(UserIDs))
	for i, v := range UserIDs {
		userIDStrs[i] = strconv.Itoa(int(v))
	}
	user_list, err := GetUsersByID(userIDStrs)
	if err != nil {
		return
	}
	users = make(map[int64]*model.User, 0)
	for _, v := range user_list {
		users[int64(v.ID)] = &model.User{
			Id:              int64(v.ID),
			Name:            v.UserName,
			FollowCount:     pointer.Int64Ptr(int64(v.FollowingCount)),
			FollowerCount:   pointer.Int64Ptr(int64(v.FollowerCount)),
			Avatar:          pointer.StringPtr(v.Avatar),
			IsFollow:        IsFollow(uint(self), v.ID),
			BackgroundImage: pointer.StringPtr(v.BackgroundImage),
			FavoriteCount:   pointer.Int64Ptr(int64(v.FavoriteCount)),
			TotalFavorited:  pointer.Int64Ptr(int64(v.TotalFavorited)),
			WorkCount:       pointer.Int64Ptr(int64(v.WorkCount)),
		}
	}
	return
}

// GetUserInfo 从缓存或数据库中获取用户信息
func GetUserInfo(userID string) (user *models.User, err error) {
	// 查询缓存中是否存在用户信息
	exist := crud.redis.HExists(context.Background(), "UserInfoCache", userID)
	if exist.Val() {
		// 从缓存中获取用户信息
		res := crud.redis.HGet(context.Background(), "UserInfoCache", userID)
		user = new(models.User)
		data, err := res.Result()
		if err != nil {
			return nil, err
		}
		// 反序列化缓存数据到用户对象
		err = sonic.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
		return user, err
	}
	// 缓存中不存在用户信息，从数据库中获取
	err = crud.mysql.Model(&models.User{}).Take(&user, userID).Error
	if err != nil {
		return nil, err
	}
	// 将用户信息添加到缓存
	err = CacheUserInfo(user)
	return
}

// IsFollow 判断用户是否关注了某个用户
func IsFollow(userID, toUserID uint) bool {
	if userID == toUserID {
		return true
	}
	ex := crud.redis.Exists(context.Background(), userRelationFollowKey(userID))
	if ex.Val() == 0 {
		var users []*models.User
		users, _ = models.GetFollowList(userID)
		CacheRelationFollows(userID, users)
		CacheUsersInfo(users)
	}
	ret := crud.redis.SIsMember(context.Background(), userRelationFollowKey(userID), toUserID).Val()
	return ret
}

// GetUsersByID 根据用户ID列表从缓存或数据库中批量获取用户信息
func GetUsersByID(userIDs []string) (users []*models.User, err error) {
	var data string
	var user *models.User
	var r redis.Cmder
	var i int
	users = make([]*models.User, len(userIDs))
	// 修复缓存清除后的bug
	if crud.redis.Exists(context.Background(), "UserInfoCache").Val() != 1 {
		var mysql_users []*models.User
		err = crud.mysql.Where("id in ?", userIDs).Find(&mysql_users).Error
		if err == nil {
			CacheUsersInfo(mysql_users)
		}
		return mysql_users, err
	}

	pipline := crud.redis.Pipeline()
	defer pipline.Close()
	// 遍历用户ID列表，批量查询缓存
	for _, id := range userIDs {
		pipline.HGet(context.Background(), "UserInfoCache", id)
	}
	// 执行缓存查询操作
	res, _ := pipline.Exec(context.Background())

	uncached_users_id := make([]string, 0)
	uncached_users_pos := make([]int, 0)
	// 遍历查询结果，进行反序列化,记录缓存中没有查询到的用户
	for i, r = range res {
		data, err = r.(*redis.StringCmd).Result()
		if err != nil {
			// 记录cache中没有查到的id
			uncached_users_id = append(uncached_users_id, userIDs[i])
			uncached_users_pos = append(uncached_users_pos, i)
			continue
		}
		user = new(models.User)
		err = sonic.Unmarshal([]byte(data), &user)
		if err != nil {
			return
		}
		users[i] = user
	}

	// 查询mysql获取缓存中没有的用户
	if len(uncached_users_id) > 0 {
		var mysql_users []*models.User
		crud.mysql.Where("id in ?", uncached_users_id).Find(&mysql_users)
		for i, v := range mysql_users {
			users[uncached_users_pos[i]] = v
		}
		// 存入缓存
		CacheUsersInfo(mysql_users)
	}
	return users, err
}

// 用户成功登录后将其信息加载到redis缓存中
func LoadUserCache(userID uint) (err error) {

	var relations []*models.FollowRelation
	// 在关系表中查询当前用户关注的对象
	crud.mysql.Where("user_id=?", userID).Preload("ToUser").Find(&relations)

	follows := make([]any, len(relations))
	followUsers := make([]any, len(relations)*2+2)
	// 遍历关系结果 使用follows被关注的用户id  使用followUsers用户信息
	for i, v := range relations {
		follows[i] = v.ToUserID
		data, _ := sonic.Marshal(v.ToUser)
		// followUsers[v.ToUserID] = string(data)
		followUsers[2*i] = v.ToUserID
		followUsers[2*i+1] = string(data)
	}
	// 获取用户自身的信息加入cache中
	var user = new(models.User)
	crud.mysql.Take(&user, userID)
	data, _ := sonic.Marshal(user)
	followUsers[len(relations)*2] = user.ID
	followUsers[len(relations)*2+1] = string(data)

	// 查询关注当前用户的用户
	crud.mysql.Where("to_user_id=?", userID).Preload("User").Find(&relations)
	followerUsers := make([]any, 2*len(relations))
	followers := make([]any, len(relations))
	for i, v := range relations {
		followers[i] = v.UserID
		if v.User.ID == 0 {
			continue
		}
		data, _ := sonic.Marshal(v.User)
		// followerUsers[v.UserID] = string(data)
		followerUsers[2*i] = v.UserID
		followerUsers[2*i+1] = string(data)
	}

	pipline := crud.redis.Pipeline()
	defer pipline.Close()
	// 加载关系表缓存
	pipline.SAdd(context.Background(), userRelationFollowKey(userID), follows...)
	pipline.SAdd(context.Background(), userRelationFollowerKey(userID), followers...)
	// 加载用户信息缓存
	pipline.HSet(context.Background(), "UserInfoCache", followUsers...)
	pipline.HSet(context.Background(), "UserInfoCache", followerUsers...)

	_, err = pipline.Exec(context.Background())
	return err
}
