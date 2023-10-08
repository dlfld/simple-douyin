package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/douyin/kitex_gen/gpt"
	"github.com/douyin/kitex_gen/gpt/chatgptservice"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"sync"
	"testing"
)

var GptService = "gpt-service"
var GptAddr = "43.130.60.218:8888"
var once sync.Once
var cli chatgptservice.Client
var err error

func NewRpcInteractionClient() (chatgptservice.Client, error) {
	once.Do(func() {
		cli, err = chatgptservice.NewClient(GptService, client.WithHostPorts(GptAddr), client.WithSuite(tracing.NewClientSuite()))
		if err != nil {
			panic(err)
		}
	})
	return cli, err
}
func TestGpt(t *testing.T) {
	cli, err := NewRpcInteractionClient()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := cli.GptChat(context.Background(), &gpt.GptChatRequest{
		Msg: "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp.Msg)
}
