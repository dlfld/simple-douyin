package models

import (
	"encoding/json"
)

// Message
//
//	@Description: 聊天消息数据模型
type Message struct {
	Id         int64  `gorm:"primarykey" json:"id"`
	FromUserID int64  `gorm:"index:idx_userid_from;not null" json:"from_user_id"`
	ToUserID   int64  `gorm:"index:idx_userid_to;not null" json:"to_user_id"`
	Content    string `gorm:"type:varchar(255);not null" json:"content"`
	CreateTime int64  `gorm:"index;not null" json:"create_time"`
}

func (m *Message) TableName() string {
	return "messages"
}

// 序列化
func (m *Message) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

// 反序列化
func (m *Message) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
