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
	"os"
	"io"
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



// // UploadVideo
// // @Description: 上传视频
// // @param bucketName 桶名字
// // @param filePath 文件路径
// // @return error
// func (minioOss *MinioService) UploadVideo(bucketName string, filePath string) error {
// 	ctx := context.Background()

// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		log.Printf("无法打开文件 %s: %v\n", filePath, err)
// 		return err
// 	}
// 	defer file.Close()

// 	fileInfo, err := file.Stat()
// 	if err != nil {
// 		log.Printf("获取不到文件中的信息 %s: %v\n", filePath, err)
// 		return err
// 	}

// 	// 创建一个新的 PutObjectOptions 结构体
// 	// opts := minio.PutObjectOptions{ContentType: "application/octet-stream"}
// 	opts := minio.PutObjectOptions{ContentType: "video/mp4"}
	
// 	// 上传文件
// 	_, err = minioOss.Client.PutObject(ctx, bucketName, fileInfo.Name(), file, fileInfo.Size(), opts)
// 	if err != nil {
// 		log.Printf("无法上传文件 %s: %v\n", filePath, err)
// 		return err
// 	}

// 	log.Printf("成功把文件 %s 上传到 %s 桶中\n", filePath, bucketName)
// 	return nil
// }

// UploadFile
// @Description: 通用的文件上传函数
// @param bucketName 桶名字
// @param filePath 文件路径
// @param contentType 文件类型(image/jpeg video/mp4)
// @return error
func (minioOss *MinioService) UploadFile(bucketName string, filePath string, contentType string) error {
	ctx := context.Background()

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Unable to open file %s: %v\n", filePath, err)
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Unable to get info of file %s: %v\n", filePath, err)
		return err
	}

	// 创建一个新的 PutObjectOptions 结构体
	opts := minio.PutObjectOptions{ContentType: contentType}

	// 上传文件
	_, err = minioOss.Client.PutObject(ctx, bucketName, fileInfo.Name(), file, fileInfo.Size(), opts)
	if err != nil {
		log.Printf("Unable to upload file %s: %v\n", filePath, err)
		return err
	}

	log.Printf("Successfully uploaded file %s to bucket %s\n", filePath, bucketName)
	return nil
}


func (minioOss *MinioService) UploadFileWithBytestream(bucketName string, reader io.Reader, fileName string, fileSize int64, contentType string) error {
	ctx := context.Background()

	// 创建一个新的 PutObjectOptions 结构体
	opts := minio.PutObjectOptions{ContentType: contentType}

	// 上传文件
	_, err := minioOss.Client.PutObject(ctx, bucketName, fileName, reader, fileSize, opts)
	if err != nil {
		log.Printf("Unable to upload file %s: %v\n", fileName, err)
		return err
	}

	log.Printf("Successfully uploaded file %s to bucket %s\n", fileName, bucketName)
	return nil
}
