package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"yuyan/internal/api"
	"yuyan/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	// Set up configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Create default config file if it doesn't exist
	if err := createDefaultConfigIfNotExists(); err != nil {
		log.Fatalf("Failed to create default config: %v", err)
	}

	// Read configuration
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Set default values
	viper.SetDefault("server.port", 8090)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.sqlite.path", "./data/yuyan.db")
}

func main() {
	// Set Gin mode
	mode := viper.GetString("server.mode")
	if mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create data directory if it doesn't exist
	if viper.GetString("database.type") == "sqlite" {
		dbPath := viper.GetString("database.sqlite.path")
		dataDir := filepath.Dir(dbPath)
		if _, err := os.Stat(dataDir); os.IsNotExist(err) {
			if err := os.MkdirAll(dataDir, 0755); err != nil {
				log.Fatalf("Failed to create data directory: %v", err)
			}
		}
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Set up static file server
	r.Static("/static", "web/static")
	r.StaticFile("/favicon.ico", "web/static/favicon.ico")

	// Define template functions
	r.SetFuncMap(template.FuncMap{
		"getCurrentLanguage": func() string {
			// This is a placeholder function to be used in templates
			// The actual language detection happens in JavaScript
			return "en"
		},
	})

	// Load HTML templates
	r.LoadHTMLGlob("web/templates/*.html")

	// Set up routes
	api.SetupRouter(r)

	// Start the server
	port := viper.GetInt("server.port")
	addr := fmt.Sprintf(":%d", port)

	log.Printf("Starting server on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// createDefaultConfigIfNotExists creates a default config file if it doesn't exist
func createDefaultConfigIfNotExists() error {
	// Check if config file exists
	if _, err := os.Stat("./config/config.yaml"); err == nil {
		return nil
	}

	// Create config directory if it doesn't exist
	if _, err := os.Stat("./config"); os.IsNotExist(err) {
		if err := os.MkdirAll("./config", 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	// Default config content
	configContent := `server:
  port: 8090
  mode: debug

database:
  type: sqlite # sqlite, mysql, postgres
  sqlite:
    path: ./data/yuyan.db
  mysql:
    host: localhost
    port: 3306
    username: root
    password: password
    dbname: yuyan
    params: charset=utf8mb4&parseTime=True&loc=Local
  postgres:
    host: localhost
    port: 5432
    username: postgres
    password: password
    dbname: yuyan
    sslmode: disable
`

	// Write default config
	if err := os.WriteFile("./config/config.yaml", []byte(strings.TrimSpace(configContent)), 0644); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	return nil
}
