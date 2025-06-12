package api

import (
	"net/http"
	"strconv"

	"yuyan/internal/database"
	"yuyan/internal/models"
	"yuyan/internal/service"

	"github.com/gin-gonic/gin"

	log "github.com/luojiego/slogx"
)

// MessageAPI handles message-related API endpoints
type MessageAPI struct {
	MessageService *service.MessageService
}

// MessageRequest represents a request to send a message
type MessageRequest struct {
	BotID   uint   `json:"bot_id" binding:"required"`
	Content string `json:"content" binding:"required"`
	Format  string `json:"format"`
}

// DashboardStats represents statistics for the dashboard
type DashboardStats struct {
	TotalBots       int64 `json:"total_bots"`
	MessagesSent    int64 `json:"messages_sent"`
	MessagesPending int64 `json:"messages_pending"`
	MessagesFailed  int64 `json:"messages_failed"`
}

// NewMessageAPI creates a new MessageAPI handler
func NewMessageAPI(messageService *service.MessageService) *MessageAPI {
	return &MessageAPI{
		MessageService: messageService,
	}
}

// GetAllMessages handles GET /api/messages
func (api *MessageAPI) GetAllMessages(c *gin.Context) {
	// Get bot type filter from query
	botType := c.Query("type")
	botIDStr := c.Query("bot_id")

	var messages interface{}
	var err error

	if botType != "" {
		messages, err = api.MessageService.GetMessagesByBotType(models.BotType(botType))
	} else if botIDStr != "" {
		botID, err := strconv.ParseUint(botIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bot ID"})
			return
		}
		messages, err = api.MessageService.GetMessagesByBotID(uint(botID))
	} else {
		messages, err = api.MessageService.GetAllMessages()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// GetRecentMessages handles GET /api/messages/recent
func (api *MessageAPI) GetRecentMessages(c *gin.Context) {
	// Get recent messages (limit to 10)
	var messages []models.Message
	result := database.DB.Preload("Bot").Order("created_at desc").Limit(10).Find(&messages)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// GetDashboardStats handles GET /api/dashboard
func (api *MessageAPI) GetDashboardStats(c *gin.Context) {
	// Initialize stats
	stats := DashboardStats{}

	// Count total active bots
	database.DB.Model(&models.Bot{}).Where("is_active = ?", true).Count(&stats.TotalBots)

	// Count messages by status
	database.DB.Model(&models.Message{}).Where("status = ?", models.MessageStatusSent).Count(&stats.MessagesSent)
	database.DB.Model(&models.Message{}).Where("status = ?", models.MessageStatusPending).Count(&stats.MessagesPending)
	database.DB.Model(&models.Message{}).Where("status = ?", models.MessageStatusFailed).Count(&stats.MessagesFailed)

	c.JSON(http.StatusOK, stats)
}

// GetMessageByID handles GET /api/messages/:id
func (api *MessageAPI) GetMessageByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	message, err := api.MessageService.GetMessageByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Format response for better frontend compatibility
	response := gin.H{
		"id":         message.ID,
		"bot_id":     message.BotID,
		"bot_name":   message.Bot.Name,
		"bot_type":   message.Bot.Type,
		"content":    message.Content,
		"format":     message.Format,
		"mentions":   message.Mentions,
		"status":     message.Status,
		"error":      message.Error,
		"sent_at":    message.SentAt,
		"created_at": message.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// SendMessage handles POST /api/messages
func (api *MessageAPI) SendMessage(c *gin.Context) {
	var req MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("Failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default to text format if not specified
	if req.Format == "" {
		req.Format = "text"
	}

	message, err := api.MessageService.SendMessage(req.BotID, req.Content, req.Format)
	if err != nil {
		log.Error("Failed to send message", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}

// DeleteMessage handles DELETE /api/messages/:id
func (api *MessageAPI) DeleteMessage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	err = api.MessageService.DeleteMessage(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}
