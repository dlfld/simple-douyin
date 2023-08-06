package crud

import (
	"fmt"
	"testing"
	"time"

	"github.com/douyin/models"
)

func TestCachedGetFollow(t *testing.T) {
	crud, err := NewCachedCRUD()
	if err != nil {
		fmt.Println(err)
	}
	// err = crud.LoadUserCache(1)
	ts := time.Now()
	u, _ := crud.RelationGetFriends(1)
	fmt.Println("crud using:", time.Since(ts))
	for _, v := range u {
		fmt.Println(v.ID)
	}
	ts = time.Now()
	models.GetFriendList(1)
	fmt.Println("sql using:", time.Since(ts))
	// ts := time.Now()

	// models.GetFollowList(1)
	// fmt.Println("sql using:", time.Since(ts))
	// ts = time.Now()
	// crud.RelationGetFollows(1)
	// fmt.Println("cache using:", time.Since(ts))
}
