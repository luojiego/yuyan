package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"yuyan/internal/models"

	"github.com/spf13/viper"
)

// TelegramBot implements a Telegram notification bot
type TelegramBot struct {
	Bot models.Bot
}

// TelegramMessage represents a message to be sent to Telegram
type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// NewTelegramBot creates a new Telegram bot instance
func NewTelegramBot(bot models.Bot) *TelegramBot {
	return &TelegramBot{
		Bot: bot,
	}
}

// Send sends a message through the Telegram bot
func (t *TelegramBot) Send(message string) error {
	// Telegram API endpoint for sending messages
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Bot.Token)

	// Create message payload
	msg := TelegramMessage{
		ChatID:    t.Bot.WebhookURL, // In Telegram case, WebhookURL field stores the chat ID
		Text:      message,
		ParseMode: "HTML", // Can be "HTML" or "Markdown"
	}

	// Convert message to JSON
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal Telegram message: %w", err)
	}

	// Create HTTP client with proxy if available
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Check if proxy is enabled in config
	if viper.GetBool("proxy.enable") {
		// Get proxy URL from config
		proxyURLStr := viper.GetString("proxy.url")
		if proxyURLStr == "" {
			// Fallback to environment variable
			proxyURLStr = os.Getenv("HTTP_PROXY")
		}

		if proxyURLStr != "" {
			// Configure proxy
			proxyURL, err := url.Parse(proxyURLStr)
			if err != nil {
				return fmt.Errorf("invalid proxy URL: %w", err)
			}

			// Set up proxy transport
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
			client.Transport = transport
		}
	}

	// Send HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(msgBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Telegram message: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("failed to send Telegram message: status %d", resp.StatusCode)
		}
		return fmt.Errorf("failed to send Telegram message: %v", result)
	}

	return nil
}
