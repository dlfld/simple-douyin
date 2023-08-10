package main

import (
	"fmt"
	"github.com/douyin/kitex_gen/message"
)

type ErrResponse struct {
	Code int32
	Err  error
}

var success = "发送成功"
var (
	ErrMysqlConn    = &ErrResponse{1001, fmt.Errorf("mysql连接错误")}
	ErrMysqlWrite   = &ErrResponse{1002, fmt.Errorf("mysql写入错误")}
	ErrMysqlRead    = &ErrResponse{1003, fmt.Errorf("mysql读取错误")}
	ErrCacheConn    = &ErrResponse{1004, fmt.Errorf("缓存连接错误")}
	ErrCacheWrite   = &ErrResponse{1005, fmt.Errorf("缓存写入错误")}
	ErrCacheRead    = &ErrResponse{1006, fmt.Errorf("缓存读取错误")}
	ErrActivateType = &ErrResponse{1007, fmt.Errorf("错误的消息发送类型")}

	ActionSuccess = &message.MessageActionResponse{0, &success}
	ChatSuccess   = &message.MessageChatResponse{0, &success, nil}
)

func HandlerActionError(err *ErrResponse) *message.MessageActionResponse {
	msg := err.Err.Error()
	return &message.MessageActionResponse{
		StatusCode: err.Code,
		StatusMsg:  &msg,
	}
}

func HandlerChatError(err *ErrResponse) *message.MessageChatResponse {
	msg := err.Err.Error()
	return &message.MessageChatResponse{
		StatusCode: err.Code,
		StatusMsg:  &msg,
	}
}
