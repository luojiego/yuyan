package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// ConfigAPI handles configuration-related API endpoints
type ConfigAPI struct{}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port int    `json:"port" binding:"required"`
	Mode string `json:"mode" binding:"required"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Type     string                 `json:"type" binding:"required"`
	SQLite   map[string]interface{} `json:"sqlite,omitempty"`
	MySQL    map[string]interface{} `json:"mysql,omitempty"`
	Postgres map[string]interface{} `json:"postgres,omitempty"`
}

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

// NewConfigAPI creates a new ConfigAPI handler
func NewConfigAPI() *ConfigAPI {
	return &ConfigAPI{}
}

// GetConfig handles GET /api/config
func (api *ConfigAPI) GetConfig(c *gin.Context) {
	config := Config{
		Server: ServerConfig{
			Port: viper.GetInt("server.port"),
			Mode: viper.GetString("server.mode"),
		},
		Database: DatabaseConfig{
			Type: viper.GetString("database.type"),
			SQLite: map[string]interface{}{
				"path": viper.GetString("database.sqlite.path"),
			},
			MySQL: map[string]interface{}{
				"host":     viper.GetString("database.mysql.host"),
				"port":     viper.GetInt("database.mysql.port"),
				"username": viper.GetString("database.mysql.username"),
				"password": viper.GetString("database.mysql.password"),
				"dbname":   viper.GetString("database.mysql.dbname"),
				"params":   viper.GetString("database.mysql.params"),
			},
			Postgres: map[string]interface{}{
				"host":     viper.GetString("database.postgres.host"),
				"port":     viper.GetInt("database.postgres.port"),
				"username": viper.GetString("database.postgres.username"),
				"password": viper.GetString("database.postgres.password"),
				"dbname":   viper.GetString("database.postgres.dbname"),
				"sslmode":  viper.GetString("database.postgres.sslmode"),
			},
		},
	}

	c.JSON(http.StatusOK, config)
}

// UpdateConfig handles PUT /api/config
func (api *ConfigAPI) UpdateConfig(c *gin.Context) {
	var req Config
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update server config
	viper.Set("server.port", req.Server.Port)
	viper.Set("server.mode", req.Server.Mode)

	// Update database config
	viper.Set("database.type", req.Database.Type)

	if req.Database.SQLite != nil {
		for key, value := range req.Database.SQLite {
			viper.Set("database.sqlite."+key, value)
		}
	}

	if req.Database.MySQL != nil {
		for key, value := range req.Database.MySQL {
			viper.Set("database.mysql."+key, value)
		}
	}

	if req.Database.Postgres != nil {
		for key, value := range req.Database.Postgres {
			viper.Set("database.postgres."+key, value)
		}
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(viper.ConfigFileUsed())
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create config directory"})
			return
		}
	}

	// Save config
	if err := viper.WriteConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully. Server restart required for changes to take effect."})
}
