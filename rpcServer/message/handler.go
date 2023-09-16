package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/douyin/common/constant"
	"github.com/douyin/kitex_gen/message"
	"github.com/douyin/kitex_gen/model"
	"github.com/go-redis/redis/v8"
	"sort"
	"strconv"
	"time"
)

// MessageServiceImpl implements the last service interface defined in the IDL.
type MessageServiceImpl struct{}

// MessageList implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageList(ctx context.Context, req *message.MessageChatRequest) (resp *message.MessageChatResponse, err error) {
	resp = message.NewMessageChatResponse()

	exists, err := bf.CheckIfUserIdExists(req.ToUserId)
	if err != nil {
		logCollector.Error(fmt.Sprintf("Message bloom_user err[%v]", err))
	} else {
		if !exists {
			constant.HandlerErr(constant.ErrBloomUser, resp)
			return resp, nil
		}
	}

	// 1. 数据库和缓存中读取消息
	// req.FromUserId是当前用户id，req.ToUserId是对方用户id
	// 得到当前用户发送给别人的消息
	sendMes := make([]*model.Message, 0)
	if req.PreMsgTime == 0 {
		sendMes, err = getSenderToReceiverMes(ctx, req.FromUserId, req.ToUserId, req.PreMsgTime)
		if err != nil {
			logCollector.Error(fmt.Sprintf("user[%d]: Failed to query the message sent by user[%d] to user[%d] where preMsgTime=%s, err=%s",
				req.ToUserId, req.FromUserId, req.ToUserId, time.Unix(req.PreMsgTime, 0).Format("2006-01-02 15:04:05"), err.Error()),
			)
			constant.HandlerErr(constant.ErrSystemBusy, resp)
			return resp, nil
		}
	}
	// 得到当前用户接收的消息
	receivedMsg, err := getSenderToReceiverMes(ctx, req.ToUserId, req.FromUserId, req.PreMsgTime)
	if err != nil {
		logCollector.Error(
			fmt.Sprintf("user[%d]: Failed to query the message sent by user[%d] to user[%d] where preMsgTime=%s, err=%s",
				req.ToUserId, req.FromUserId, req.ToUserId, time.Unix(req.PreMsgTime, 0).Format("2006-01-02 15:04:05"), err.Error(),
			))
		constant.HandlerErr(constant.ErrSystemBusy, resp)
		return resp, nil
	}
	// 冷启动，为了不重复显示消息，只返回时间
	if len(sendMes) == 1 && len(receivedMsg) == 0 {
		sendMes[0].Content = ""
	}
	messageList := append(sendMes, receivedMsg...)
	sort.Slice(messageList, func(i, j int) bool {
		return messageList[i].CreateTime < messageList[j].CreateTime
	})
	// 2. 得到返回值
	logCollector.Info(fmt.Sprintf("user[%d] successfully queried %d messages sent from user[%d] b since %s",
		req.ToUserId, len(messageList), req.FromUserId, time.Unix(req.PreMsgTime, 0).Format("2006-01-02 15:04:05")),
	)
	resp.MessageList = messageList
	return resp, nil
}

// getSenderToReceiverMes 得到sender发送给receiver的消息
func getSenderToReceiverMes(ctx context.Context, sender, receiver, preMsgTime int64) ([]*model.Message, error) {
	// 1. 统计缓存中数据量
	keyFrom := fmt.Sprintf("%s:%d:%d", messageCacheTable, sender, receiver)
	var num int64
	if err := cache.getRecordNum(ctx, keyFrom, &num).Err(); err != nil {
		return nil, err
	}
	res := make([]*model.Message, 0)
	// 2. 读取数据
	switch num {
	// 2.1 缓存中无数据，直接去mysql数据库中读取
	case 0:
		if err := db.Table((&model.Message{}).TableName()).
			Where("from_user_id=? and to_user_id=? and create_time > ?", sender, receiver, preMsgTime).
			Scan(&res).Error; err != nil {
			return nil, err
		}
	// 2.2 缓存中有数据，读缓存
	default:
		// 得到缓存中最小的时间分数
		minScore, err := cache.getMinScore(ctx, keyFrom)
		if err != nil {
			return nil, err
		}
		mysqlRecords := make([]*model.Message, 0)
		// 缓存数据无法cover全部要查询的数据，则需要去mysql查询没cover的那部分数据
		if minScore > preMsgTime {
			// 这里的req.ToUserId是对方用户id是指消息来自哪里，对应数据库消息记录的from_user_id
			if err = db.Table((&model.Message{}).TableName()).
				Where("from_user_id=? and to_user_id=? and create_time > ? and create_time < ? ",
					sender, receiver, preMsgTime, minScore).
				Scan(&mysqlRecords).Error; err != nil {
				return nil, err
			}
		}
		// 读取缓存中记录
		cacheRecordsFrom := cache.ZRangeByScore(ctx, keyFrom, &redis.ZRangeBy{
			Min: strconv.FormatInt(preMsgTime+1, 10),
			Max: "inf",
		})
		if cache.Err() != nil {
			return nil, err
		}
		res = append(mysqlRecords, cacheRecordsFrom...)
	}
	return res, nil
}

// SendMessage implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) SendMessage(ctx context.Context, req *message.MessageActionRequest) (resp *message.MessageActionResponse, err error) {
	// 1. 消息合法类型验证
	resp = message.NewMessageActionResponse()

	exists, err := bf.CheckIfUserIdExists(req.ToUserId)
	if err != nil {
		logCollector.Error(fmt.Sprintf("Message bloom_user err[%v]", err))
	} else {
		if !exists {
			constant.HandlerErr(constant.ErrBloomUser, resp)
			return resp, nil
		}
	}

	if req.ActionType != 1 {
		constant.HandlerErr(constant.ErrUnsupportedOperation, resp)
		return resp, nil
	}
	// 2. 转为存入数据库的消息记录，并添加时间戳
	messageRecord := model.Message{
		FromUserId: req.FromUserId,
		ToUserId:   req.ToUserId,
		Content:    req.Content,
		CreateTime: time.Now().UnixMilli(),
	}
	tx := db.Begin()
	// 3. 存入mysql数据库表
	if err = tx.Create(&messageRecord).Error; err != nil {
		tx.Rollback()
		logCollector.Error(
			fmt.Sprintf("failed to create message record in mysql, err=%s", err.Error()),
		)
		constant.HandlerErr(constant.ErrSystemBusy, resp)
		return resp, nil
	}
	logCollector.Info(fmt.Sprintf("The message from user[%d] to user[%d] was successfully stored in MySQL, content: %s",
		messageRecord.FromUserId, messageRecord.ToUserId, messageRecord.Content))
	// 4. 序列化后存入缓存
	key := fmt.Sprintf("%s:%d:%d", messageCacheTable, req.FromUserId, req.ToUserId)
	record, err := json.Marshal(messageRecord)
	if err != nil {
		tx.Rollback()
		logCollector.Error(fmt.Sprintf("failed to marshal, err: %s", err.Error()))
		constant.HandlerErr(constant.ErrSystemBusy, resp)
		return resp, nil
	}
	if err = cache.ZAdd(ctx, key, float64(messageRecord.CreateTime), record, time.Hour*24*7).Err(); err != nil {
		tx.Rollback()
		logCollector.Error(fmt.Sprintf("failed to create cache, err: %s", err.Error()))
		constant.HandlerErr(constant.ErrSystemBusy, resp)
		return nil, err
	}
	tx.Commit()
	logCollector.Info(fmt.Sprintf("The message from user[%d] to user[%d] was successfully stored in Redis, key: %s, content: %s",
		messageRecord.FromUserId, messageRecord.ToUserId, key, messageRecord.Content))
	// 5. 缓存消息记录超过设定的阈值，就从缓存中删除5天之前的数据
	go cache.keepDataNum(ctx, key)
	return resp, nil
}
