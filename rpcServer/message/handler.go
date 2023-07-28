package main

import (
	"context"
	"fmt"
	"github.com/douyin/kitex_gen/message"
)

// MessageServiceImpl implements the last service interface defined in the IDL.
type MessageServiceImpl struct{}

// MessageList implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageList(ctx context.Context, req *message.MessageChatRequest) (resp *message.MessageChatResponse, err error) {
	// TODO: Your code here...
	if req.Token == "right" {
		resp.StatusCode = 8888
	}
	return
}

// SendMessage implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) SendMessage(ctx context.Context, req *message.MessageActionRequest) (resp *message.MessageActionResponse, err error) {
	// TODO: Your code here...
	msg := "发送消息成功"
	if req.Token == "right" {
		resp = &message.MessageActionResponse{
			StatusMsg: &msg,
		}
	}
	fmt.Println("1111")
	return
}
