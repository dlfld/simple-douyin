package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
)

func DiscoverService(serviceName string) ([]string, error) {
	// 创建etcd客户端连接
	cli, err := GetEtcdCli()
	defer cli.Close()

	// 查询服务信息
	key := fmt.Sprintf("/services/%s", serviceName)
	resp, err := cli.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	// 解析服务地址
	addresses := make([]string, 0)
	for _, kv := range resp.Kvs {
		addresses = append(addresses, string(kv.Value))
	}

	return addresses, nil
}
