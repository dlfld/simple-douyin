// Package cosService @author:戴林峰
// @date:2023/8/25
// @node:
package cosService

import (
	"context"
	"errors"
	"github.com/douyin/common/conf"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

var once sync.Once

// CosService Cos Service 对象
type CosService struct {
	Client *cos.Client
}

// GetCos
//
//	@Description: 获取cos对象
//	@return *cos
//	@return error
func GetCos() (*CosService, error) {
	cosCase := &CosService{}
	client, err := cosCase.GetClient()
	cosCase.Client = client.(*cos.Client)
	if err != nil {
		return nil, err
	}
	return cosCase, nil
}

// GetClient 获取ossclient 单例模式
func (cosOss *CosService) GetClient() (interface{}, error) {
	var client *cos.Client
	var err error

	once.Do(func() {
		u, _ := url.Parse(conf.CosConfig.Url)
		su, _ := url.Parse(conf.CosConfig.ReginUrl)
		b := &cos.BaseURL{BucketURL: u, ServiceURL: su}
		client = cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  os.Getenv(conf.CosConfig.SecretID),  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
				SecretKey: os.Getenv(conf.CosConfig.SecretKey), // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
			},
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
func (cosOss *CosService) CreateBucket(bucketName string) error {
	if len(bucketName) <= 0 {
		return errors.New("bucketName invalid")
	}
	return nil
}

// DeleteBucket
//
//	@Description: 删除桶
//	@param bucketName 桶名字
//	@return error
func (cosOss *CosService) DeleteBucket(bucketName string) error {
	return nil
}

// UploadFile
// @Description: 通用的文件上传函数
// @param bucketName 桶名字
// @param filePath 文件路径
// @param contentType 文件类型(image/jpeg video/mp4)
// @return error
func (cosOss *CosService) UploadFile(bucketName string, filePath string, contentType string) error {
	return nil
}

func (cosOss *CosService) UploadFileWithBytestream(bucketName string, reader io.Reader, fileName string, fileSize int64, contentType string) error {
	_, err := cosOss.Client.Object.Put(context.Background(), fileName, reader, nil)
	if err != nil {
		log.Printf("Unable to upload file %s: %v\n", fileName, err)
		return err
	}

	log.Printf("Successfully uploaded file %s to bucket %s\n", fileName, bucketName)
	return nil
}

// RemoveObject
// @Description: 删除对象
// @receiver cosOss
// @param bulkName 桶名
// @param objectName 对象名
// @return error
func (cosOss *CosService) RemoveObject(bulkName, objectName string) error {
	return nil
}

//	GetPlayUrl
//
// @Description:	获取播放链接
// @param key
// @return string
// @return error
func (cosOss *CosService) GetPlayUrl(key string) (string, error) {
	url := cosOss.Client.Object.GetObjectURL(key)
	return url.String(), nil
}
