package service

import (
	"fmt"

	"yuyan/internal/database"
	"yuyan/internal/models"
)

// BotService provides methods for managing bot configurations
type BotService struct{}

// GetAllBots retrieves all bot configurations
func (s *BotService) GetAllBots() ([]models.Bot, error) {
	var bots []models.Bot
	result := database.DB.Find(&bots)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get bots: %w", result.Error)
	}
	return bots, nil
}

// GetBotsByType retrieves bot configurations by type
func (s *BotService) GetBotsByType(botType models.BotType) ([]models.Bot, error) {
	var bots []models.Bot
	result := database.DB.Where("type = ?", botType).Find(&bots)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get bots by type: %w", result.Error)
	}
	return bots, nil
}

// GetBotByID retrieves a bot configuration by ID
func (s *BotService) GetBotByID(id uint) (*models.Bot, error) {
	var bot models.Bot
	result := database.DB.First(&bot, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get bot by ID: %w", result.Error)
	}
	return &bot, nil
}

// CreateBot creates a new bot configuration
func (s *BotService) CreateBot(bot models.Bot) (*models.Bot, error) {
	result := database.DB.Create(&bot)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create bot: %w", result.Error)
	}
	return &bot, nil
}

// UpdateBot updates an existing bot configuration
func (s *BotService) UpdateBot(bot models.Bot) (*models.Bot, error) {
	result := database.DB.Save(&bot)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update bot: %w", result.Error)
	}
	return &bot, nil
}

// DeleteBot deletes a bot configuration
func (s *BotService) DeleteBot(id uint) error {
	result := database.DB.Delete(&models.Bot{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete bot: %w", result.Error)
	}
	return nil
}
