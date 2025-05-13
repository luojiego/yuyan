package database

import (
	"fmt"
	"log"

	"yuyan/internal/models"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() error {
	var err error
	dbType := viper.GetString("database.type")

	switch dbType {
	case "sqlite":
		path := viper.GetString("database.sqlite.path")
		DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to connect to SQLite database: %w", err)
		}
		log.Println("Connected to SQLite database")

	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
			viper.GetString("database.mysql.username"),
			viper.GetString("database.mysql.password"),
			viper.GetString("database.mysql.host"),
			viper.GetInt("database.mysql.port"),
			viper.GetString("database.mysql.dbname"),
			viper.GetString("database.mysql.params"),
		)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to connect to MySQL database: %w", err)
		}
		log.Println("Connected to MySQL database")

	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			viper.GetString("database.postgres.host"),
			viper.GetInt("database.postgres.port"),
			viper.GetString("database.postgres.username"),
			viper.GetString("database.postgres.password"),
			viper.GetString("database.postgres.dbname"),
			viper.GetString("database.postgres.sslmode"),
		)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
		}
		log.Println("Connected to PostgreSQL database")

	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Auto migrate models
	err = DB.AutoMigrate(&models.Bot{}, &models.Message{})
	if err != nil {
		return fmt.Errorf("failed to auto-migrate models: %w", err)
	}

	return nil
}
