package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"yuyan/internal/models"

	"github.com/spf13/viper"
)

// TelegramBot implements a Telegram notification bot
type TelegramBot struct {
	Bot models.Bot
}

// TelegramMessage represents a message to be sent to Telegram
type TelegramMessage struct {
	ChatID                string `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

// NewTelegramBot creates a new Telegram bot instance
func NewTelegramBot(bot models.Bot) *TelegramBot {
	return &TelegramBot{
		Bot: bot,
	}
}

// validateMessage checks if the message is valid
func (t *TelegramBot) validateMessage(message string) error {
	if !utf8.ValidString(message) {
		return fmt.Errorf("invalid UTF-8 encoding in message")
	}

	if len(message) > 4096 {
		return fmt.Errorf("message exceeds Telegram's 4096 character limit")
	}

	return nil
}

// parseTelegramMentions extracts @mention information from a message and formats for Telegram
// Returns a formatted message with proper HTML tags for mentions
func parseTelegramMentions(message string) string {
	// Check for @all mention and add announcement at the beginning if present
	if strings.Contains(message, "@all") {
		message = "<b>📢 Attention Everyone! 📢</b><br>" + message
	}

	// Regular expression to find @mentions with usernames
	// Telegram usernames are 5-32 characters and can contain a-z, A-Z, 0-9 and underscores
	re := regexp.MustCompile(`@([a-zA-Z0-9_]{5,32})`)
	message = re.ReplaceAllString(message, `<a href="https://t.me/$1">@$1</a>`)

	// For phone numbers (@12345678901), convert to a special format
	// since Telegram doesn't allow mentioning by phone number directly
	rePhone := regexp.MustCompile(`@(\d{10,15})`)
	message = rePhone.ReplaceAllString(message, `<b>@$1</b>`)

	return message
}

// Send sends a message through the Telegram bot
func (t *TelegramBot) Send(message string) error {
	// Validate message
	if err := t.validateMessage(message); err != nil {
		return fmt.Errorf("message validation failed: %w", err)
	}

	// Process mentions in the message
	formattedMessage := parseTelegramMentions(message)

	// Create message payload
	msg := TelegramMessage{
		ChatID:                t.Bot.WebhookURL,
		Text:                  formattedMessage,
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
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
		proxyURLStr := viper.GetString("proxy.url")
		if proxyURLStr == "" {
			proxyURLStr = os.Getenv("HTTP_PROXY")
		}

		if proxyURLStr != "" {
			proxyURL, err := url.Parse(proxyURLStr)
			if err != nil {
				return fmt.Errorf("invalid proxy URL: %w", err)
			}

			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
			client.Transport = transport
		}
	}

	// Send HTTP request
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Bot.Token)
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
