package etcd

import (
	"go.etcd.io/etcd/clientv3"
	"sync"
	"time"
)

var cli *clientv3.Client
var once sync.Once
var err error

// GetEtcdCli // 创建etcd客户端连接
func GetEtcdCli() (*clientv3.Client, error) {
	once.Do(func() {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   []string{"http://localhost:2379"}, // etcd服务器地址
			DialTimeout: 5 * time.Second,
		})
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}
