package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/douyin/common/constant"
	"github.com/douyin/kitex_gen/message"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

// MessageServiceImpl implements the last service interface defined in the IDL.
type MessageServiceImpl struct{}

// MessageList implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageList(ctx context.Context, req *message.MessageChatRequest) (resp *message.MessageChatResponse, err error) {
	resp = message.NewMessageChatResponse()
	messageList := make([]*model.Message, 0)
	// 1. 统计缓存中数据量
	key := fmt.Sprintf("%s:%d:%d", messageCacheTable, req.ToUserId, req.FromUserId)
	var num int64
	if err = cache.getRecordNum(ctx, key, &num).Err(); err != nil {
		logger.Errorf("failed count cache number, err: %s", err.Error())
		return nil, err
	}
	dbMessageList := make([]*models.Message, 0)
	// 2. 读取数据
	switch num {
	// 2.1 缓存中无数据，直接去mysql数据库中读取
	case 0:
		if err = db.Table((&models.Message{}).TableName()).
			Where("from_user_id=? and to_user_id=? and created_time >= ?", req.ToUserId, req.FromUserId, req.PreMsgTime).
			Scan(&dbMessageList).Error; err != nil {
			logger.Errorf("query data when created_time >= %d failed, err: %s", req.PreMsgTime, err.Error())
			return nil, err
		}
	// 2.2 缓存中有数据，读缓存
	default:
		// 得到缓存中最小的时间分数
		minScore, err := cache.getMinScore(ctx, key)
		if err != nil {
			logger.Errorf("get minimal score failed, err: %s", err.Error())
			return nil, err
		}
		// 缓存数据无法cover全部要查询的数据，则需要去mysql查询没cover的那部分数据
		mysqlRecords := make([]*models.Message, 0)
		if minScore > req.PreMsgTime {
			// 这里的req.ToUserId是对方用户id是指消息来自哪里，对应数据库消息记录的from_user_id
			if err = db.Table((&models.Message{}).TableName()).
				Where("from_user_id=? and to_user_id=? and created_time between ? and ? ", req.ToUserId, req.FromUserId, req.PreMsgTime, minScore-1).
				Scan(&mysqlRecords).Error; err != nil {
				return nil, err
			}
		}
		// 读取缓存中记录
		cacheRecords := cache.ZRangeByScore(ctx, key, &redis.ZRangeBy{
			Min: strconv.FormatInt(req.PreMsgTime, 10),
			Max: "inf",
		})
		if cache.Err() != nil {
			logger.Errorf("failed to query message list, err: %s", err.Error())
			return nil, err
		}
		dbMessageList = append(dbMessageList, mysqlRecords...)
		dbMessageList = append(dbMessageList, cacheRecords...)
	}
	// 3. 解析得到返回值
	for i := 0; i < len(dbMessageList); i++ {
		strTime := time.Unix(dbMessageList[i].CreatedTime, 0).Format("2006-01-02 15:04:05")
		messageList = append(messageList, &model.Message{
			Id:         dbMessageList[i].Id,
			ToUserId:   dbMessageList[i].ToUserID,
			FromUserId: dbMessageList[i].FromUserID,
			Content:    dbMessageList[i].Content,
			CreateTime: &strTime,
		})
	}
	logger.Infof("user[%d] successfully queried %d messages sent from user[%d] b since %s",
		req.ToUserId, len(messageList), req.FromUserId, time.Unix(req.PreMsgTime, 0).Format("2006-01-02 15:04:05"))
	resp.MessageList = messageList
	return resp, nil
}

// SendMessage implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) SendMessage(ctx context.Context, req *message.MessageActionRequest) (resp *message.MessageActionResponse, err error) {
	// 1. 消息合法类型验证
	resp = message.NewMessageActionResponse()
	if req.ActionType != 1 {
		return resp, constant.ErrUnsupportedOperation
	}
	// 2. 转为存入数据库的消息记录，并添加时间戳
	messageRecord := models.Message{
		FromUserID:  req.FromUserId,
		ToUserID:    req.ToUserId,
		Content:     req.Content,
		CreatedTime: time.Now().Unix(),
	}
	// 3. 存入mysql数据库表
	if err = db.Create(&messageRecord).Error; err != nil {
		logger.Errorf("failed to create message record in mysql, err: %s", err.Error())
		return resp, constant.ErrSystemBusy
	}
	logger.Infof("The message from user[%d] to user[%d] was successfully stored in MySQL, content: %s",
		messageRecord.FromUserID, messageRecord.ToUserID, messageRecord.Content)
	// 4. 序列化后存入缓存
	key := fmt.Sprintf("%s:%d:%d", messageCacheTable, req.ToUserId, req.FromUserId)
	record, err := json.Marshal(messageRecord)
	if err != nil {
		logger.Errorf("failed to marshal, err: %s", err.Error())
		return nil, constant.ErrSystemBusy
	}
	if err = cache.ZAdd(ctx, key, float64(messageRecord.CreatedTime), record, time.Hour*24*7).Err(); err != nil {
		logger.Errorf("create cache failed, err: %s", err.Error())
		return nil, err
	}
	logger.Infof("The message from user[%d] to user[%s] was successfully stored in Redis, key: %s, content: %s",
		messageRecord.FromUserID, messageRecord.ToUserID, key, messageRecord.Content)
	// 5. 缓存消息记录超过设定的阈值，就从缓存中删除5天之前的数据
	go cache.keepDataNum(ctx, key)
	return resp, nil
}
