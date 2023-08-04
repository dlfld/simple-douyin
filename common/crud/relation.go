package crud

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/douyin/models"
	"github.com/go-redis/redis/v8"
)

func userRelationFollowKey(userID uint) string {
	return fmt.Sprintf("relation_follow_%d", userID)
}
func userRelationFollowerKey(userID uint) string {
	return fmt.Sprintf("relation_follower_%d", userID)
}

// RelationFollow 在缓存中建立关注关系
func (c *CachedCRUD) RelationFollow(userID, toUserID uint) (err error) {
	ex := c.redis.HExists(context.Background(), "UserInfoCache", strconv.Itoa(int(toUserID)))
	if !ex.Val() {
		var user models.User
		c.mysql.Take(&user, toUserID)
		c.CacheUserInfo(&user)
	}
	pipline := c.redis.Pipeline()
	defer pipline.Close()
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
func (c *CachedCRUD) RelationUnFollow(userID, toUserID uint) (err error) {
	pipline := c.redis.Pipeline()
	defer pipline.Close()

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
func (c *CachedCRUD) RelationGetFollows(userID uint) (userList []*models.User, err error) {
	// 获取用户关注列表
	res := c.redis.SMembers(context.Background(), userRelationFollowKey(userID))
	if res.Err() != nil {
		return nil, res.Err()
	}

	// 获取关注用户信息
	var ids = res.Val()
	userList, err = c.GetUsersByID(ids)
	return userList, err
}

// RelationGetFollowers 获取用户的粉丝列表
func (c *CachedCRUD) RelationGetFollowers(userID uint) (userList []*models.User, err error) {
	// 获取用户的粉丝列表
	res := c.redis.SMembers(context.Background(), userRelationFollowerKey(userID))
	if res.Err() != nil {
		return nil, res.Err()
	}
	var ids = res.Val()
	// 获取粉丝用户信息
	userList, err = c.GetUsersByID(ids)
	return userList, err
}

// RelationGetFriends 获取用户的好友列表
func (c *CachedCRUD) RelationGetFriends(userID uint) (userList []*models.User, err error) {
	// 获取交集，即用户的好友列表
	res := c.redis.SInter(context.Background(), userRelationFollowerKey(userID), userRelationFollowKey(userID))
	if res.Err() != nil {
		return nil, res.Err()
	}
	var ids = res.Val()
	// 获取好友用户信息
	userList, err = c.GetUsersByID(ids)
	return userList, err
}

// CacheUserInfo 将用户信息存入缓存
func (c *CachedCRUD) CacheUserInfo(user *models.User) (err error) {
	// 序列化用户信息
	data, err := sonic.Marshal(user)
	if err != nil {
		return err
	}
	// 将用户信息存入缓存的哈希表中
	res := c.redis.HSet(context.Background(), "UserInfoCache", fmt.Sprintf("%d", user.ID), string(data))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

// GetUserInfo 从缓存或数据库中获取用户信息
func (c *CachedCRUD) GetUserInfo(userID string) (user *models.User, err error) {
	// 查询缓存中是否存在用户信息
	exist := c.redis.HExists(context.Background(), "UserInfoCache", userID)
	if exist.Val() {
		// 从缓存中获取用户信息
		res := c.redis.HGet(context.Background(), "UserInfoCache", userID)
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
	err = c.mysql.Model(&models.User{}).Take(&user, userID).Error
	if err != nil {
		return nil, err
	}
	// 将用户信息添加到缓存
	err = c.CacheUserInfo(user)
	return
}

// GetUsersByID 根据用户ID列表从缓存或数据库中批量获取用户信息
func (c *CachedCRUD) GetUsersByID(userIDs []string) (users []*models.User, err error) {
	var data string
	var user *models.User
	var r redis.Cmder
	var i int
	users = make([]*models.User, len(userIDs))
	pipline := c.redis.Pipeline()
	defer pipline.Close()
	// 遍历用户ID列表，批量查询缓存
	for _, id := range userIDs {
		pipline.HGet(context.Background(), "UserInfoCache", id)
	}
	// 执行缓存查询操作
	res, err := pipline.Exec(context.Background())
	if err != nil {
		// fmt.Println("%v\n", err.Error())
		return
	}

	// 遍历查询结果，进行反序列化
	for i, r = range res {
		data, err = r.(*redis.StringCmd).Result()
		if err != nil {
			return nil, err
		}
		user = new(models.User)
		err = sonic.Unmarshal([]byte(data), &user)
		if err != nil {
			return
		}
		users[i] = user
	}
	return users, err
}

// 用户成功登录后将其信息加载到redis缓存中
func (c *CachedCRUD) LoadUserCache(userID uint) (err error) {

	var relations []*models.FollowRelation
	// 在关系表中查询当前用户关注的对象
	c.mysql.Where("user_id=?", userID).Preload("ToUser").Find(&relations)

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
	c.mysql.Take(&user, userID)
	data, _ := sonic.Marshal(user)
	followUsers[len(relations)*2] = user.ID
	followUsers[len(relations)*2+1] = string(data)

	// 查询关注当前用户的用户
	c.mysql.Where("to_user_id=?", userID).Preload("User").Find(&relations)
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

	pipline := c.redis.Pipeline()
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
