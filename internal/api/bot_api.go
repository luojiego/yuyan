package api

import (
	"net/http"
	"net/url"
	"strconv"

	"yuyan/internal/models"
	"yuyan/internal/service"

	"github.com/gin-gonic/gin"
)

// BotAPI handles bot-related API endpoints
type BotAPI struct {
	BotService *service.BotService
}

// BotRequest represents a request to create or update a bot
type BotRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Token       string `json:"token" binding:"required"`
	Secret      string `json:"secret"`
	WebhookURL  string `json:"webhook_url"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// NewBotAPI creates a new BotAPI handler
func NewBotAPI(botService *service.BotService) *BotAPI {
	return &BotAPI{
		BotService: botService,
	}
}

// GetAllBots handles GET /api/bots
func (api *BotAPI) GetAllBots(c *gin.Context) {
	// Get bot type filter from query
	botType := c.Query("type")

	var bots []models.Bot
	var err error

	if botType != "" {
		bots, err = api.BotService.GetBotsByType(models.BotType(botType))
	} else {
		bots, err = api.BotService.GetAllBots()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bots)
}

// GetBotByID handles GET /api/bots/:id
func (api *BotAPI) GetBotByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bot ID"})
		return
	}

	bot, err := api.BotService.GetBotByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bot not found"})
		return
	}

	c.JSON(http.StatusOK, bot)
}

// CreateBot handles POST /api/bots
func (api *BotAPI) CreateBot(c *gin.Context) {
	var req BotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate webhook URL if provided
	if req.WebhookURL != "" {
		if _, err := url.Parse(req.WebhookURL); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook URL: " + err.Error()})
			return
		}
	}

	bot := models.Bot{
		Name:        req.Name,
		Type:        models.BotType(req.Type),
		Token:       req.Token,
		Secret:      req.Secret,
		WebhookURL:  req.WebhookURL,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	created, err := api.BotService.CreateBot(bot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// UpdateBot handles PUT /api/bots/:id
func (api *BotAPI) UpdateBot(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bot ID"})
		return
	}

	var req BotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate webhook URL if provided
	if req.WebhookURL != "" {
		if _, err := url.Parse(req.WebhookURL); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook URL: " + err.Error()})
			return
		}
	}

	// Get existing bot
	existing, err := api.BotService.GetBotByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bot not found"})
		return
	}

	// Update fields
	existing.Name = req.Name
	existing.Type = models.BotType(req.Type)
	existing.Token = req.Token
	existing.Secret = req.Secret
	existing.WebhookURL = req.WebhookURL
	existing.Description = req.Description
	existing.IsActive = req.IsActive

	updated, err := api.BotService.UpdateBot(*existing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteBot handles DELETE /api/bots/:id
func (api *BotAPI) DeleteBot(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bot ID"})
		return
	}

	err = api.BotService.DeleteBot(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bot deleted successfully"})
}
