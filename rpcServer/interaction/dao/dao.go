package dao

import (
	"github.com/douyin/common/bloom"
)

type Dao struct {
	Mysql       *mysql
	Redis       *redis
	BloomFilter *bloom.Filter
}

func NewDao() (dao *Dao) {
	dao = &Dao{
		Mysql:       NewMysql(),
		Redis:       NewRedis(),
		BloomFilter: bloom.NewBloom(),
	}
	return dao
}
