package etcd

import (
	"fmt"
	"testing"

	"github.com/douyin/common/conf"
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
