package models

import (
	"time"

	"gorm.io/gorm"
)

// Message
//
//	@Description: 聊天消息数据模型
type Message struct {
	ID         uint      `gorm:"primarykey"`
	CreatedAt  time.Time `gorm:"index;not null" json:"create_time"`
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	FromUser   User           `gorm:"foreignkey:FromUserID;" json:"from_user,omitempty"`
	FromUserID uint           `gorm:"index:idx_userid_from;not null" json:"from_user_id"`
	ToUser     User           `gorm:"foreignkey:ToUserID;" json:"to_user,omitempty"`
	ToUserID   uint           `gorm:"index:idx_userid_from;index:idx_userid_to;not null" json:"to_user_id"`
	Content    string         `gorm:"type:varchar(255);not null" json:"content"`
}

func (Message) TableName() string {
	return "messages"
}
