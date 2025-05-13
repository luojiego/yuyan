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

// parseTelegramMentions extracts @mention information from a message and formats for Telegram
// Returns a formatted message with proper HTML tags for mentions
func parseTelegramMentions(message string) string {
	// Check for @all mention (Telegram doesn't have a direct equivalent,
	// we'll replace with a general announcement message)
	if strings.Contains(message, "@all") {
		message = strings.Replace(message, "@all", "<b>📢 Attention Everyone! 📢</b>", -1)
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
	// Telegram API endpoint for sending messages
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Bot.Token)

	// Process mentions in the message
	formattedMessage := parseTelegramMentions(message)

	// Create message payload
	msg := TelegramMessage{
		ChatID:                t.Bot.WebhookURL, // In Telegram case, WebhookURL field stores the chat ID
		Text:                  formattedMessage,
		ParseMode:             "HTML", // Use HTML parse mode to support formatted messages
		DisableWebPagePreview: true,   // Prevent link previews for user mentions
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
