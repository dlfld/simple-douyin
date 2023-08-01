// @author:戴林峰
// @date:2023/8/1
// @node:
package oss

import (
	"fmt"
	"testing"
)

func TestService_CreateBucket(t *testing.T) {
	ossService, _ := GetOssService()
	_ = ossService.CreateBucket("test1")
}
func TestGetOssService(t *testing.T) {
	service, _ := GetOssService()
	fmt.Printf("%+v", service)
}

func TestService_DeleteBucket(t *testing.T) {
	service, _ := GetOssService()
	_ = service.DeleteBucket("test1")
}
