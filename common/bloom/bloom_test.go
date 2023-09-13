package bloom

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestBloomFilter(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "",
		DB:       0,
	})

	key := "my_bloom_filter"
	element := "?ASdsa"

	result, err := rdb.Do(context.Background(), "BF.ADD", key, element).Bool()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	result, err = rdb.Do(context.Background(), "BF.EXISTS", key, element).Bool()
	if err != nil {
		fmt.Println("检查元素失败:", err)
	} else {
		fmt.Println("元素是否存在:", result)
	}
}

func TestBloomFilter2(t *testing.T) {
	bloom := NewBloom()
	ids := []int64{1, 3, 5, 7, 9}
	bloom.AddVideoIds(ids)
	bloom.AddCommentIds(ids)
	bloom.AddUserIds(ids)
	for i := 0; i < 10; i++ {
		exists, err := bloom.CheckIfVideoIdExists(int64(i))
		if err != nil {
			fmt.Println("id : ", i, " err : ", err)
		}
		fmt.Println(fmt.Sprintf("id(%d) 是否存在(%v)", i, exists))
	}

	for i := 0; i < 10; i++ {
		exists, err := bloom.CheckIfCommentIdExists(int64(i))
		if err != nil {
			fmt.Println("id : ", i, " err : ", err)
		}
		fmt.Println(fmt.Sprintf("id(%d) 是否存在(%v)", i, exists))
	}

	for i := 0; i < 10; i++ {
		exists, err := bloom.CheckIfUserIdExists(int64(i))
		if err != nil {
			fmt.Println("id : ", i, " err : ", err)
		}
		fmt.Println(fmt.Sprintf("id(%d) 是否存在(%v)", i, exists))
	}
}
