package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/douyin/common/mysql"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
)

func TestORM(t *testing.T) {
	// var user models.User
	// var user2 models.User
	// var relations models.FollowRelation
	db, _ := mysql.NewMysqlConn()
	// db.Take(&user, 1)
	// db.Take(&user2, 3)
	// var followers []*models.FollowRelation
	// relation := models.FollowRelation{
	// 	UserID:   user.ID,
	// 	ToUserID: user2.ID,
	// }
	// db.Create(&relation)
	var followers models.FollowRelation
	db.Preload("User").Preload("ToUser").First(&followers)
	fmt.Println(followers.ToUser)
	fmt.Println(followers.User)

}

func TestGetAllFollows(t *testing.T) {
	db, _ := mysql.NewMysqlConn()
	userID := 1
	var follows []*models.FollowRelation
	db.Where("user_id=?", userID).Preload("ToUser").Find(&follows)
	fmt.Println(follows)
	for _, v := range follows {
		fmt.Println("follos:", v.ToUserID, v.ToUser.UserName)
	}

}

func TestGetAllFollowers(t *testing.T) {
	db, _ := mysql.NewMysqlConn()
	userID := 1
	var follows []*models.FollowRelation
	db.Where("to_user_id=?", userID).Preload("User").Find(&follows)
	fmt.Println(follows)
	for _, v := range follows {
		fmt.Println("followers:", v.UserID, v.ToUser.UserName)
	}

}

func TestFollow(t *testing.T) {
	db, _ := mysql.NewMysqlConn()
	var toUserId, userID uint
	userID = 7
	toUserId = 1
	relation := models.FollowRelation{
		UserID:   userID,
		ToUserID: toUserId,
	}

	// t.Log("relation id:", relation.ID)
	db.Where("to_user_id=? and user_id = ?", toUserId, userID).FirstOrCreate(&relation)
	// t.Log("relation id:", relation.ID)

}
func TestUserExist(t *testing.T) {
	db, _ := mysql.NewMysqlConn()
	var user models.User
	var n int64
	db.Find(&user, 1).Count(&n)

	fmt.Println(n)
	fmt.Println(user)
}

func TestGetFriends(t *testing.T) {
	db, _ := mysql.NewMysqlConn()
	var relations []*models.FollowRelation
	var userList []*models.User
	// userID := 1
	// select * from relations where user_id=1 and to_user_id in (select user_id from relations where to_user_id=1);
	subQuery := db.Model(&models.FollowRelation{}).Select("user_id").Where("to_user_id = ?", 1)
	db.Where("user_id = ? AND to_user_id IN (?)", 1, subQuery).Preload("ToUser").Find(&relations)
	for _, v := range relations {
		userList = append(userList, &v.ToUser)
		fmt.Println(v.ToUser.UserName, v.ToUser.ID)
	}
	fmt.Println(userList)
}

func TestConvert(t *testing.T) {
	userList, _ := models.GetFriendList(1)
	data, _ := json.Marshal(userList)
	var kitexList []*model.User
	fmt.Println(string(data))
	json.Unmarshal(data, &kitexList)
	fmt.Println(kitexList[0])

}

// func TestVideo(t *testing.T) {
// 	db, _ := mysql.NewMysqlConn()
// 	videos := make([]*models.Video, 0)
// 	// db.Find(&videos, 1)
// 	db.Where("id = ?", "1").Find(&videos)
// 	fmt.Println("video:", videos[0])
// }

// func TestStructConvert(t *testing.T) {
// 	type GormStruct struct {
// 		Id     int `json:"idaaa"`
// 		Name   string
// 		Gender bool
// 	}
// 	type KitexStruct struct {
// 		ID       int
// 		Nameaa   string
// 		Genderbb bool `json:"dddd"`
// 	}

// 	gormstruct := GormStruct{
// 		Id:     1,
// 		Name:   "hello",
// 		Gender: true,
// 	}

// 	kitestruct := (*KitexStruct)(unsafe.Pointer(&gormstruct))

// 	field, _ := reflect.TypeOf(kitestruct).Elem().FieldByName("Genderbb")
// 	fmt.Printf("field: %#v tag: %#v\n", field.Name, field.Tag)
// 	fmt.Printf("type: %T value: %#v\n", kitestruct, kitestruct)

// }

// func BenchmarkPointerConvert(b *testing.B) {
// 	type GormStruct struct {
// 		Id     int `json:"idaaa"`
// 		Name   string
// 		Gender bool
// 	}
// 	type KitexStruct struct {
// 		ID       int
// 		Nameaa   string
// 		Genderbb bool `json:"dddd"`
// 	}
// 	gormstruct := GormStruct{
// 		Id:     1,
// 		Name:   "hello",
// 		Gender: true,
// 	}
// 	b.ResetTimer()
// 	for i:=0;i<b.N;i++{
// 		kitestruct:=(*KitexStruct)(unsafe.Pointer(&gormstruct))

// 		kitestruct.ID=2;
// 	}

// }
