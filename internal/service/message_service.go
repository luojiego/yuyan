package service

import (
	"fmt"
	"regexp"
	"strings"
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

// extractMentions parses @mentions from a message and returns them as a comma-separated string
func extractMentions(content string) string {
	re := regexp.MustCompile(`@([\w\d]+)`)
	matches := re.FindAllStringSubmatch(content, -1)

	mentions := make([]string, 0)

	for _, match := range matches {
		if len(match) > 1 && match[1] != "all" {
			mentions = append(mentions, match[1])
		}
	}

	return strings.Join(mentions, ",")
}

// SendMessage sends a message through the specified bot
func (s *MessageService) SendMessage(botID uint, content string, format string) (*models.Message, error) {
	// Get bot configuration
	botConfig, err := s.BotService.GetBotByID(botID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot: %w", err)
	}

	// Extract mentions from the content
	mentions := extractMentions(content)

	// Create message record
	message := models.Message{
		BotID:    botID,
		Content:  content,
		Status:   models.MessageStatusProcessing,
		Format:   format,
		Mentions: mentions,
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

// DeleteMessage deletes a message by ID
func (s *MessageService) DeleteMessage(id uint) error {
	// Check if message exists
	var message models.Message
	result := database.DB.First(&message, id)
	if result.Error != nil {
		return fmt.Errorf("failed to find message: %w", result.Error)
	}

	// Delete the message (soft delete with gorm)
	result = database.DB.Delete(&message)
	if result.Error != nil {
		return fmt.Errorf("failed to delete message: %w", result.Error)
	}

	return nil
}
