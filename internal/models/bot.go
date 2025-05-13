package models

import (
	"time"

	"gorm.io/gorm"
)

// BotType represents the type of notification bot
type BotType string

const (
	BotTypeDingDing BotType = "dingtalk"
	BotTypeTelegram BotType = "telegram"
)

// Bot represents a notification bot configuration
type Bot struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Type        BotType        `json:"type" gorm:"size:50;not null"`
	Token       string         `json:"token" gorm:"size:255;not null"`
	Secret      string         `json:"secret" gorm:"size:255"`
	WebhookURL  string         `json:"webhook_url" gorm:"size:255"`
	Description string         `json:"description" gorm:"size:255"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
