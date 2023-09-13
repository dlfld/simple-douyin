package etcd

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/douyin/common/conf"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestDiscoverService(t *testing.T) {
	a := DiscoverService(conf.UserService.Name)
	fmt.Println(a)
}

func TestRegisterService(t *testing.T) {
	RegisterService(conf.InteractionService.Name, conf.InteractionService.Addr)
	a := DiscoverService(conf.InteractionService.Name)
	fmt.Println(a)
}

func TestGetAll(t *testing.T) {
	// 创建etcd客户端连接
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"}, // etcd服务器地址
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	key := fmt.Sprintf("/douyin/")

	resp, err := cli.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, kv := range resp.Kvs {
		fmt.Println(string(kv.Key), string(kv.Value))
		cli.Delete(context.Background(), string(kv.Key))
	}

}

func TestGetIP(t *testing.T) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		return
	}

	ips, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, ip := range ips {
		// 将服务信息注册到etcd中
		fmt.Println(ip.String())
	}
}
