package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/douyin/common/conf"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func DiscoverService(serviceName string) (addr []string) {
	// 创建etcd客户端连接
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{conf.EtcdConfig.Addr}, // etcd服务器地址
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
