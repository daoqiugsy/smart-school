package model

import (
	"time"
)

// Notification 通知信息
type Notification struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"size:100;not null" json:"title"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Type       int       `gorm:"not null" json:"type"`               // 0:系统通知 1:课程通知 2:考试通知 3:作业通知 4:行政通知
	Priority   int       `gorm:"not null;default:1" json:"priority"` // 0:低 1:中 2:高
	SourceID   uint      `gorm:"default:null" json:"source_id"`      // 来源ID，如课程ID、考试ID等
	SourceType string    `gorm:"size:50" json:"source_type"`         // 来源类型，如Course、Exam等
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Notification) TableName() string {
	return "notifications"
}

// UserNotification 用户通知关联
type UserNotification struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	UserID         uint         `gorm:"not null;index" json:"user_id"`
	NotificationID uint         `gorm:"not null;index" json:"notification_id"`
	IsRead         bool         `gorm:"not null;default:false" json:"is_read"`     // 是否已读
	ReadTime       time.Time    `json:"read_time"`                                 // 阅读时间
	DeliveryStatus int          `gorm:"not null;default:0" json:"delivery_status"` // 0:待发送 1:已发送 2:发送失败
	DeliveryType   int          `gorm:"not null" json:"delivery_type"`             // 0:APP 1:短信 2:邮箱 3:语音
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	User           User         `gorm:"foreignKey:UserID" json:"user"`
	Notification   Notification `gorm:"foreignKey:NotificationID" json:"notification"`
}

// TableName 指定表名
func (UserNotification) TableName() string {
	return "user_notifications"
}
