package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func DiscoverService(serviceName string) (addr []string) {
	// 创建etcd客户端连接
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://42.192.46.30:2379"}, // etcd服务器地址
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// 查询服务信息
	key := fmt.Sprintf("/douyin/%s", serviceName)
	resp, err := cli.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}

	// 解析服务地址
	addresses := make([]string, 0)
	for _, kv := range resp.Kvs {
		addresses = append(addresses, string(kv.Value))
	}

	//// 客户端负载均衡
	//rand.Seed(1)
	//randomIndex := rand.Intn(len(addresses))
	//randomAddress := addresses[randomIndex]

	return addresses
}
