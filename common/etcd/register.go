package etcd

import (
	"context"
	"fmt"

	"go.etcd.io/etcd/clientv3"
)

func RegisterService(serviceName string, serviceAddress string) error {
	// 创建etcd客户端连接
	cli, err := GetEtcdCli()
	defer cli.Close()

	// 创建租约
	resp, err := cli.Grant(context.Background(), 5)
	if err != nil {
		return err
	}

	// 将服务信息注册到etcd中
	key := fmt.Sprintf("/services/%s/%s", serviceName, serviceAddress)
	_, err = cli.Put(context.Background(), key, serviceAddress, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	// 定期刷新租约
	keepAliveCh, err := cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
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

	return nil
}
