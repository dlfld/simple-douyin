package handler

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"io"
//	"net"
//	"net/http"
//	"strconv"
//	"sync"
//	"sync/atomic"
//	"time"
//)
//
//var chatConnMap = sync.Map{}
//
//func RunMessageServer() {
//	listen, err := net.Listen("tcp", "127.0.0.1:9090")
//	if err != nil {
//		fmt.Printf("Run message sever failed: %v\n", err)
//		return
//	}
//
//	for {
//		conn, err := listen.Accept()
//		if err != nil {
//			fmt.Printf("Accept conn failed: %v\n", err)
//			continue
//		}
//
//		go process(conn)
//	}
//}
//
//func process(conn net.Conn) {
//	defer conn.Close()
//
//	var buf [256]byte
//	for {
//		n, err := conn.Read(buf[:])
//		if n == 0 {
//			if err == io.EOF {
//				break
//			}
//			fmt.Printf("Read message failed: %v\n", err)
//			continue
//		}
//
//		var event = MessageSendEvent{}
//		_ = json.Unmarshal(buf[:n], &event)
//		fmt.Printf("Receive Message：%+v\n", event)
//
//		fromChatKey := fmt.Sprintf("%d_%d", event.UserId, event.ToUserId)
//		if len(event.MsgContent) == 0 {
//			chatConnMap.Store(fromChatKey, conn)
//			continue
//		}
//
//		toChatKey := fmt.Sprintf("%d_%d", event.ToUserId, event.UserId)
//		writeConn, exist := chatConnMap.Load(toChatKey)
//		if !exist {
//			fmt.Printf("User %d offline\n", event.ToUserId)
//			continue
//		}
//
//		pushEvent := MessagePushEvent{
//			FromUserId: event.UserId,
//			MsgContent: event.MsgContent,
//		}
//		pushData, _ := json.Marshal(pushEvent)
//		_, err = writeConn.(net.Conn).Write(pushData)
//		if err != nil {
//			fmt.Printf("Push message failed: %v\n", err)
//		}
//	}
//}
//
//var tempChat = map[string][]Message{}
//
//var messageIdSequence = int64(1)
//
//type ChatResponse struct {
//	Response
//	MessageList []Message `json:"message_list"`
//}
//
//// MessageAction no practical effect, just check if token is valid
//func MessageAction(c *gin.Context) {
//	toUserId := c.Query("to_user_id")
//	content := c.Query("content")
//	//userId := int64(c.GetUint("userID"))
//	//toUserId := c.PostForm("to_user_id")
//	//content := c.PostForm("content")
//	userId := int64(c.GetUint("userID"))
//	userIdB, _ := strconv.Atoi(toUserId)
//	chatKey := genChatKey(userId, int64(userIdB))
//
//	atomic.AddInt64(&messageIdSequence, 1)
//	curMessage := Message{
//		Id:         messageIdSequence,
//		Content:    content,
//		CreateTime: time.Now().Unix(),
//	}
//
//	if messages, exist := tempChat[chatKey]; exist {
//		tempChat[chatKey] = append(messages, curMessage)
//	} else {
//		tempChat[chatKey] = []Message{curMessage}
//	}
//	fmt.Println(chatKey, ": ", tempChat[chatKey])
//	c.JSON(http.StatusOK, Response{StatusCode: 0})
//}
//
//// MessageChat all users have same follow list
//func MessageChat(c *gin.Context) {
//	toUserId := c.Query("to_user_id")
//	userId := int64(c.GetUint("userID"))
//	userIdB, _ := strconv.Atoi(toUserId)
//	chatKey := genChatKey(userId, int64(userIdB))
//	fmt.Println(tempChat[chatKey])
//	c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: tempChat[chatKey]})
//}
//
//func genChatKey(userIdA int64, userIdB int64) string {
//	if userIdA > userIdB {
//		return fmt.Sprintf("%d_%d", userIdB, userIdA)
//	}
//	return fmt.Sprintf("%d_%d", userIdA, userIdB)
//}
//
//type Response struct {
//	StatusCode int32  `json:"status_code"`
//	StatusMsg  string `json:"status_msg,omitempty"`
//}
//
//type Video struct {
//	Id            int64  `json:"id,omitempty"`
//	Author        User   `json:"author"`
//	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
//	CoverUrl      string `json:"cover_url,omitempty"`
//	FavoriteCount int64  `json:"favorite_count,omitempty"`
//	CommentCount  int64  `json:"comment_count,omitempty"`
//	IsFavorite    bool   `json:"is_favorite,omitempty"`
//}
//
//type Comment struct {
//	Id         int64  `json:"id,omitempty"`
//	User       User   `json:"user"`
//	Content    string `json:"content,omitempty"`
//	CreateDate string `json:"create_date,omitempty"`
//}
//
//type User struct {
//	Id            int64  `json:"id,omitempty"`
//	Name          string `json:"name,omitempty"`
//	FollowCount   int64  `json:"follow_count,omitempty"`
//	FollowerCount int64  `json:"follower_count,omitempty"`
//	IsFollow      bool   `json:"is_follow,omitempty"`
//}
//
//type Message struct {
//	Id         int64  `json:"id,omitempty"`
//	Content    string `json:"content,omitempty"`
//	CreateTime int64  `json:"create_time,omitempty"`
//}
//
//type MessageSendEvent struct {
//	UserId     int64  `json:"user_id,omitempty"`
//	ToUserId   int64  `json:"to_user_id,omitempty"`
//	MsgContent string `json:"msg_content,omitempty"`
//}
//
//type MessagePushEvent struct {
//	FromUserId int64  `json:"user_id,omitempty"`
//	MsgContent string `json:"msg_content,omitempty"`
//}
