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
		api.GET("/messages/:id", messageAPI.GetMessageByID)
		api.POST("/messages", messageAPI.SendMessage)

		// Config endpoints
		api.GET("/config", configAPI.GetConfig)
		api.PUT("/config", configAPI.UpdateConfig)
	}

	// Frontend routes
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")

	// Define frontend routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Notification Service",
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
