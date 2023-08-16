/*
*

	@author:孟令亚
	@date:2023/8/9
	@node:

*
*/
package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	"github.com/douyin/kitex_gen/user/userservice"
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", conf.UserService.Addr)
	svr := userservice.NewServer(new(UserServiceImpl), server.WithServiceAddr(addr))
	if err != nil {
		log.Println(err.Error())
	}
	err = svr.Run()
	//if err != nil {
	//	log.Println(err.Error())
	//}
}
