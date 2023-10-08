/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"net"

	"github.com/douyin/idl/pb"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8888, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

var client = openai.NewClient("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")

// SayHello implements helloworld.GreeterServer
func (s *server) GptChat(ctx context.Context, in *pb.GptRequest) (*pb.GptReply, error) {
	gptResp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: in.Msg,
				},
			},
		},
	)
	var resp *pb.GptReply
	if err != nil {
		response(&resp, ErrChatCode, err.Error())
		fmt.Println("err: ", err.Error())
		return resp, nil
	}
	fmt.Println("return: ", gptResp.Choices[0].Message.Content)
	response(&resp, SuccessCode, gptResp.Choices[0].Message.Content)
	return resp, nil
}

type statusCode int32

const (
	SuccessCode statusCode = 0 + iota
	ErrNullContent
	ErrChatCode
)

func response(resp **pb.GptReply, status statusCode, msg string) {
	switch status {
	case SuccessCode:
		fmt.Println("1")
		*resp = &pb.GptReply{
			Msg:    msg,
			Status: int32(status),
		}
	case ErrNullContent:
		*resp = &pb.GptReply{
			Msg:    "null content",
			Status: int32(status),
		}
	case ErrChatCode:
		*resp = &pb.GptReply{
			Msg:    fmt.Sprintf("ChatCompletion error: %s\n", msg),
			Status: int32(status),
		}
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
