// Package minioService @author:戴林峰
// @date:2023/7/30
// @node:

package minioService

import (
	"context"
	"errors"
	"fmt"
	"github.com/douyin/common/conf"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"sync"
)

var once sync.Once

// MinioService Minio Service 对象
type MinioService struct {
	Client *minio.Client
}

// GetMinio
//
//	@Description: 获取minio对象
//	@return *Minio
//	@return error
func GetMinio() (*MinioService, error) {
	minioCase := &MinioService{}
	client, err := minioCase.GetClient()
	minioCase.Client = client.(*minio.Client)
	if err != nil {
		return nil, err
	}
	return minioCase, nil
}

// GetClient 获取ossclient 单例模式
func (minioOss *MinioService) GetClient() (interface{}, error) {
	var client *minio.Client
	var err error
	once.Do(func() {
		client, err = minio.New(conf.MinioConfig.EndPoint, &minio.Options{
			Creds:  credentials.NewStaticV4(conf.MinioConfig.AccessKeyId, conf.MinioConfig.SecretAccessKey, ""),
			Secure: conf.MinioConfig.UseSSL,
		})
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

// CreateBucket
//
//	@Description: 创建桶
//	@param bucketName 桶名字
//	@return error
func (minioOss *MinioService) CreateBucket(bucketName string) error {
	if len(bucketName) <= 0 {
		return errors.New("bucketName invalid")
	}
	ctx := context.Background()
	if err := minioOss.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
		exist, err := minioOss.Client.BucketExists(ctx, bucketName)
		if err != nil {
			log.Printf("%+v\n", err)
			return err
		}
		if exist {
			log.Printf("Bucket %s 已经存在！", bucketName)
			return nil
		}
	}
	log.Printf("创建Bucket %s 成功\n", bucketName)
	return nil
}

// DeleteBucket
//
//	@Description: 删除桶
//	@param bucketName 桶名字
//	@return error
func (minioOss *MinioService) DeleteBucket(bucketName string) error {
	err := minioOss.Client.RemoveBucket(context.Background(), bucketName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
