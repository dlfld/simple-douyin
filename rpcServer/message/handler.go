package main

import (
	"context"
	"github.com/douyin/kitex_gen/message"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
	"time"
)

// MessageServiceImpl implements the last service interface defined in the IDL.
type MessageServiceImpl struct{}

// MessageList implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageList(ctx context.Context, req *message.MessageChatRequest) (resp *message.MessageChatResponse, err error) {
	// TODO: Your code here...
	messageList := make([]*model.Message, 0)
	db.Table(models.Message{}.TableName()).Where("from_user_id=? and to_user_id=?", req.User.UserId, req.ToUserId).Find(&messageList)
	if db.Error != nil {
		return HandlerChatError(ErrMysqlRead), db.Error
	}
	resp = ChatSuccess
	resp.MessageList = messageList
	return resp, nil
}

// SendMessage implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) SendMessage(ctx context.Context, req *message.MessageActionRequest) (resp *message.MessageActionResponse, err error) {
	// TODO: Your code here...
	if req.ActionType != 1 {
		return HandlerActionError(ErrActivateType), nil
	}
	curTime := time.Now().UnixNano()
	messageRecord := models.Message{
		FromUserID: *req.User.UserId,
		ToUserID:   req.ToUserId,
		Content:    req.Content,
		CreatedAt:  &curTime,
	}
	if err = db.Table(messageRecord.TableName()).Create(&messageRecord).Error; err != nil {
		return HandlerActionError(ErrMysqlWrite), err
	}
	resp = ActionSuccess
	return resp, nil
}
