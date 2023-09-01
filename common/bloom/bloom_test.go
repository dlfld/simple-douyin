package bloom

import (
	"fmt"
	"github.com/douyin/common/mysql"
	"testing"
)

func TestBloomFilter(t *testing.T) {
	c, _ := mysql.NewMysqlConn()
	filter := NewBloomFilter()
	filter.PreLoadAll(c)
	Ids := []int64{1, 3, 5, 6, 77, 33}
	for _, v := range Ids {
		fmt.Println(filter.IfVideoIdExists(v))
	}
	fmt.Println("-------------------")
	filter.AddCommentIds(77)
	for _, v := range Ids {
		fmt.Println(filter.IfCommentIdExists(v))
	}
	fmt.Println("-------------------")
}
