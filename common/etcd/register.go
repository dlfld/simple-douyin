package etcd

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/douyin/common/conf"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const grantTime = 3600 * 24 * 100

func RegisterService(serviceName string, serviceAddress string) {
	// 创建etcd客户端连接
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{conf.EtcdConfig.Addr}, // etcd服务器地址
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// 创建租约
	resp, err := cli.Grant(context.Background(), grantTime)
	if err != nil {
		panic(err)
	}
	// 获取当前服务器ip
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
		newip := strings.ReplaceAll(serviceAddress, "0.0.0.0", ip.String())
		key := fmt.Sprintf("/douyin/%s/%s", serviceName, newip)
		_, err = cli.Put(context.Background(), key, newip, clientv3.WithLease(resp.ID))
		if err != nil {
			panic(err)
		}
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
