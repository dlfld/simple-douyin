package main

import (
	"context"
	"fmt"
	"github.com/douyin/kitex_gen/gpt"
	"github.com/sashabaranov/go-openai"
)

// ChatgptServiceImpl implements the last service interface defined in the IDL.
type ChatgptServiceImpl struct{}

var client = openai.NewClient("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")

// GptChat implements the ChatgptServiceImpl interface.
func (s *ChatgptServiceImpl) GptChat(ctx context.Context, req *gpt.GptChatRequest) (resp *gpt.GptChatResponse, err error) {
	// TODO: Your code here...
	gptResp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: req.Msg,
				},
			},
		},
	)
	if err != nil {
		response(resp, ErrChatCode, err.Error())
		return
	}
	response(resp, SuccessCode, gptResp.Choices[0].Message.Content)
	return
}

type statusCode int32

const (
	SuccessCode statusCode = 0 + iota
	ErrNullContent
	ErrChatCode
)

func response(resp *gpt.GptChatResponse, status statusCode, msg string) {
	switch status {
	case SuccessCode:
		resp = &gpt.GptChatResponse{
			Msg:    msg,
			Status: int32(status),
		}
	case ErrNullContent:
		resp = &gpt.GptChatResponse{
			Msg:    "null content",
			Status: int32(status),
		}
	case ErrChatCode:
		resp = &gpt.GptChatResponse{
			Msg:    fmt.Sprintf("ChatCompletion error: %s\n", msg),
			Status: int32(status),
		}
	}
}
