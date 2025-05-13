package api

import (
	"net/http"
	"strconv"

	"yuyan/internal/models"
	"yuyan/internal/service"

	"github.com/gin-gonic/gin"
)

// MessageAPI handles message-related API endpoints
type MessageAPI struct {
	MessageService *service.MessageService
}

// MessageRequest represents a request to send a message
type MessageRequest struct {
	BotID   uint   `json:"bot_id" binding:"required"`
	Content string `json:"content" binding:"required"`
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

	c.JSON(http.StatusOK, message)
}

// SendMessage handles POST /api/messages
func (api *MessageAPI) SendMessage(c *gin.Context) {
	var req MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := api.MessageService.SendMessage(req.BotID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}
