package service

import (
	"fmt"
	"time"

	"yuyan/internal/bot"
	"yuyan/internal/database"
	"yuyan/internal/models"
)

// MessageService provides methods for managing messages
type MessageService struct {
	BotService *BotService
}

// GetAllMessages retrieves all messages
func (s *MessageService) GetAllMessages() ([]models.Message, error) {
	var messages []models.Message
	result := database.DB.Preload("Bot").Find(&messages)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get messages: %w", result.Error)
	}
	return messages, nil
}

// GetMessagesByBotType retrieves messages by bot type
func (s *MessageService) GetMessagesByBotType(botType models.BotType) ([]models.Message, error) {
	var messages []models.Message
	result := database.DB.Preload("Bot").
		Joins("JOIN bots ON messages.bot_id = bots.id").
		Where("bots.type = ?", botType).
		Find(&messages)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get messages by bot type: %w", result.Error)
	}
	return messages, nil
}

// GetMessagesByBotID retrieves messages by bot ID
func (s *MessageService) GetMessagesByBotID(botID uint) ([]models.Message, error) {
	var messages []models.Message
	result := database.DB.Preload("Bot").Where("bot_id = ?", botID).Find(&messages)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get messages by bot ID: %w", result.Error)
	}
	return messages, nil
}

// GetMessageByID retrieves a message by ID
func (s *MessageService) GetMessageByID(id uint) (*models.Message, error) {
	var message models.Message
	result := database.DB.Preload("Bot").First(&message, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get message by ID: %w", result.Error)
	}
	return &message, nil
}

// SendMessage sends a message through the specified bot
func (s *MessageService) SendMessage(botID uint, content string) (*models.Message, error) {
	// Get bot configuration
	botConfig, err := s.BotService.GetBotByID(botID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot: %w", err)
	}

	// Create message record
	message := models.Message{
		BotID:   botID,
		Content: content,
		Status:  models.MessageStatusProcessing,
	}

	// Save message to database
	result := database.DB.Create(&message)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create message: %w", result.Error)
	}

	// Create bot instance
	botInstance, err := bot.NewBot(*botConfig)
	if err != nil {
		// Update message status to failed
		message.Status = models.MessageStatusFailed
		message.Error = err.Error()
		database.DB.Save(&message)
		return &message, fmt.Errorf("failed to create bot instance: %w", err)
	}

	// Send message
	err = botInstance.Send(content)
	now := time.Now()
	message.SentAt = &now

	if err != nil {
		// Update message status to failed
		message.Status = models.MessageStatusFailed
		message.Error = err.Error()
		database.DB.Save(&message)
		return &message, fmt.Errorf("failed to send message: %w", err)
	}

	// Update message status to sent
	message.Status = models.MessageStatusSent
	database.DB.Save(&message)

	return &message, nil
}
