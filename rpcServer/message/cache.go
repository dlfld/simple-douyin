package main

import (
	"context"
	"github.com/douyin/kitex_gen/model"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type Cache struct {
	rdb *redis.Client
	err error
}

// ZAdd 添加数据放入缓存
func (ca *Cache) ZAdd(ctx context.Context, key string, score float64, member interface{}, expireTime time.Duration) *Cache {
	if err := ca.rdb.ZAdd(ctx, key, &redis.Z{Score: score, Member: member}).Err(); err != nil {
		ca.err = err
		return ca
	}
	if err := ca.rdb.Expire(ctx, key, expireTime).Err(); err != nil {
		ca.err = err
	}
	return ca
}

// ZRangeByScore 根据分数从缓存中查询记录
func (ca *Cache) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) []*model.Message {
	res := make([]*model.Message, 0)
	if err := ca.rdb.ZRangeByScore(ctx, key, opt).ScanSlice(&res); err != nil {
		ca.err = err
		return nil
	}
	return res
}

// 数量超过阈值就删除距离当前时间最长的25%左右的数据
func (ca *Cache) keepDataNum(ctx context.Context, key string) *Cache {
	num, err := ca.rdb.ZCard(ctx, key).Result()
	if err != nil {
		ca.err = err
		return ca
	}
	if num > maxCacheMessageNum {
		// 默认按时间递增排序，移除下标0开始到maxCacheMessageNum/4的元素
		records := make([]*model.Message, 0)
		if err = ca.rdb.ZRange(ctx, key, maxCacheMessageNum/2, maxCacheMessageNum/2).ScanSlice(&records); err != nil {
			ca.err = err
			return ca
		}
		if err = ca.rdb.ZRemRangeByScore(ctx, key, "-inf", strconv.FormatInt(records[0].CreateTime, 10)).Err(); err != nil {
			ca.err = err
		}
	}
	return ca

}

// 得到缓存中最小的分数
func (ca *Cache) getMinScore(ctx context.Context, key string) (int64, error) {
	records := make([]*model.Message, 0)
	err := ca.rdb.ZRange(ctx, key, 0, 0).ScanSlice(&records)
	if err != nil {
		return 0, err
	}
	return records[0].CreateTime, nil
}

// 得到缓存中记录的数量
func (ca *Cache) getRecordNum(ctx context.Context, key string, res *int64) *Cache {
	num, err := ca.rdb.ZCard(ctx, key).Result()
	if err != nil {
		ca.err = err
		return ca
	}
	*res = num
	return ca
}

func (ca *Cache) Err() error {
	return ca.err
}
