// @author:戴林峰
// @date:2023/9/2
// @node:
package crud

import (
	"context"
	"fmt"
	"testing"
)

func TestCacheRelationFollowers(t *testing.T) {
	result, err := crud.redis.HGetAll(context.Background(), "asd").Result()
	if len(result) == 0 {
		print(true)
	}
	fmt.Printf("%+v,%+v", result, err)

}
