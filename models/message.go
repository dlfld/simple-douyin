package models

// Message
//
//	@Description: 聊天消息数据模型
type Message struct {
	Id         int64  `gorm:"primarykey"`
	ToUserID   int64  `gorm:"index:idx_userid_to;not null" json:"to_user_id"`
	FromUserID int64  `gorm:"index:idx_userid_from;not null" json:"from_user_id"`
	Content    string `gorm:"type:varchar(255);not null" json:"content"`
	CreatedAt  *int64 `gorm:"index;not null" json:"create_time"`
}

func (Message) TableName() string {
	return "messages"
}
