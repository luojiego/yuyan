package models

import (
	"time"

	"gorm.io/gorm"
)

// MessageStatus represents the status of a notification message
type MessageStatus string

const (
	MessageStatusPending    MessageStatus = "pending"
	MessageStatusSent       MessageStatus = "sent"
	MessageStatusFailed     MessageStatus = "failed"
	MessageStatusProcessing MessageStatus = "processing"
)

// Message represents a notification message
type Message struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	BotID     uint           `json:"bot_id" gorm:"index;not null"`
	Bot       Bot            `json:"bot" gorm:"foreignKey:BotID"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Format    string         `json:"format" gorm:"size:20;default:'text'"`
	Mentions  string         `json:"mentions" gorm:"type:text"` // Stores @mentions (mobile numbers or userIds)
	Status    MessageStatus  `json:"status" gorm:"size:20;default:'pending'"`
	Error     string         `json:"error" gorm:"size:255"`
	SentAt    *time.Time     `json:"sent_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
