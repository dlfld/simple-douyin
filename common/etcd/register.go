package etcd

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func RegisterService(serviceName string, serviceAddress string) {
	// 创建etcd客户端连接
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"}, // etcd服务器地址
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// 创建租约
	resp, err := cli.Grant(context.Background(), 5)
	if err != nil {
		panic(err)
	}

	// 将服务信息注册到etcd中
	key := fmt.Sprintf("/services/%s/%s", serviceName, serviceAddress)
	_, err = cli.Put(context.Background(), key, serviceAddress, clientv3.WithLease(resp.ID))
	if err != nil {
		panic(err)
	}

	// 定期刷新租约
	keepAliveCh, err := cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case _, ok := <-keepAliveCh:
				if !ok {
					return
				}
			}
		}
	}()

	return
}
