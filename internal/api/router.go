package api

import (
	"yuyan/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the API routes
func SetupRouter(r *gin.Engine) {
	// Create services
	botService := &service.BotService{}
	messageService := &service.MessageService{
		BotService: botService,
	}

	// Create API handlers
	botAPI := NewBotAPI(botService)
	messageAPI := NewMessageAPI(messageService)
	configAPI := NewConfigAPI()

	// API routes
	api := r.Group("/api")
	{
		// Bot endpoints
		api.GET("/bots", botAPI.GetAllBots)
		api.GET("/bots/:id", botAPI.GetBotByID)
		api.POST("/bots", botAPI.CreateBot)
		api.PUT("/bots/:id", botAPI.UpdateBot)
		api.DELETE("/bots/:id", botAPI.DeleteBot)

		// Message endpoints
		api.GET("/messages", messageAPI.GetAllMessages)
		api.GET("/messages/recent", messageAPI.GetRecentMessages)
		api.GET("/messages/:id", messageAPI.GetMessageByID)
		api.POST("/messages", messageAPI.SendMessage)
		api.DELETE("/messages/:id", messageAPI.DeleteMessage)

		// Dashboard endpoint
		api.GET("/dashboard", messageAPI.GetDashboardStats)

		// Config endpoints
		api.GET("/config", configAPI.GetConfig)
		api.PUT("/config", configAPI.UpdateConfig)
	}

	// Serve static files
	// r.Static("/static", "./web/static")  // Removed to avoid conflict with main.go

	// Web routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Dashboard",
		})
	})

	r.GET("/bots", func(c *gin.Context) {
		c.HTML(200, "bots.html", gin.H{
			"title": "Bot Management",
		})
	})

	r.GET("/messages", func(c *gin.Context) {
		c.HTML(200, "messages.html", gin.H{
			"title": "Message History",
		})
	})

	r.GET("/settings", func(c *gin.Context) {
		c.HTML(200, "settings.html", gin.H{
			"title": "System Settings",
		})
	})
}
