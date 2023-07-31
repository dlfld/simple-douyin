// Package minio @author:戴林峰
// @date:2023/7/30
// @node:
package main

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

var minioClient *minio.Client
var once sync.Once

// init 初始化minio的连接对象
func init() {
	client, err := getClient()
	if err != nil {
		log.Fatal("minio客户端初始化失败！")
	}
	minioClient = client
}

// getClient
//
//	@Description: 创建minio链接客户端
//	@return *minio.Client
//	@return error
func getClient() (*minio.Client, error) {
	// Initialize minio client object.
	var err error
	once.Do(func() {
		minioClient, err = minio.New(conf.MinioConfig.EndPoint, &minio.Options{
			Creds:  credentials.NewStaticV4(conf.MinioConfig.AccessKeyId, conf.MinioConfig.SecretAccessKey, ""),
			Secure: conf.MinioConfig.UseSSL,
		})
	})

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return minioClient, nil
}

// CreateBucket
//
//	@Description: 创建桶
//	@param bucketName 桶名字
//	@return error
func CreateBucket(bucketName string) error {
	if len(bucketName) <= 0 {
		return errors.New("bucketName invalid")
	}
	ctx := context.Background()
	if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
		exist, err := minioClient.BucketExists(ctx, bucketName)
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

func DeleteBucket(bucketName string) error {
	err := minioClient.RemoveBucket(context.Background(), bucketName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
