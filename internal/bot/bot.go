package bot

import (
	"fmt"

	"yuyan/internal/models"
)

// Bot defines the interface for all notification bots
type Bot interface {
	Send(message string) error
}

// NewBot creates a new bot instance based on bot type
func NewBot(botConfig models.Bot) (Bot, error) {
	switch botConfig.Type {
	case models.BotTypeDingDing:
		return NewDingTalkBot(botConfig), nil
	case models.BotTypeTelegram:
		return NewTelegramBot(botConfig), nil
	default:
		return nil, fmt.Errorf("unsupported bot type: %s", botConfig.Type)
	}
}
