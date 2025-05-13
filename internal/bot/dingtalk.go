package bot

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"yuyan/internal/models"
)

// DingTalkBot implements a DingTalk notification bot
type DingTalkBot struct {
	Bot models.Bot
}

// DingTalkMessage represents the message structure for DingTalk API
type DingTalkMessage struct {
	MsgType  string                 `json:"msgtype"`
	Text     map[string]string      `json:"text,omitempty"`
	Markdown map[string]string      `json:"markdown,omitempty"`
	At       map[string]interface{} `json:"at,omitempty"`
}

// DingTalkResponse represents the response structure from DingTalk API
type DingTalkResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// NewDingTalkBot creates a new DingTalk bot instance
func NewDingTalkBot(bot models.Bot) *DingTalkBot {
	return &DingTalkBot{
		Bot: bot,
	}
}

// parseMentions extracts @mention information from a message
// Returns:
// - processedMessage: the message with @mentions properly formatted
// - mobiles: list of mobile numbers to be mentioned
// - isAtAll: whether @all was mentioned
func parseMentions(message string) (string, []string, bool) {
	// Check for @all mention
	isAtAll := strings.Contains(message, "@all")

	// Regular expression to find @mentions
	// Matches patterns like @13800138000 or @user
	re := regexp.MustCompile(`@([\w\d]+)`)
	matches := re.FindAllStringSubmatch(message, -1)

	mobiles := make([]string, 0)

	// Process each match
	for _, match := range matches {
		if len(match) > 1 {
			mention := match[1]
			// Skip "all" as it's handled separately
			if mention == "all" {
				continue
			}

			// If the mention looks like a phone number (all digits), add it to mobiles
			if regexp.MustCompile(`^\d+$`).MatchString(mention) {
				mobiles = append(mobiles, mention)
			}
			// Note: For userIds, we would need to handle them differently
			// but DingTalk primarily uses phone numbers for mentions
		}
	}

	return message, mobiles, isAtAll
}

// Send sends a message through the DingTalk bot
func (d *DingTalkBot) Send(message string) error {
	// Validate webhook URL
	if d.Bot.WebhookURL == "" {
		return fmt.Errorf("webhook URL is required for DingTalk bot")
	}

	// Validate webhook URL format
	webhookURL := d.Bot.WebhookURL
	if _, err := url.Parse(webhookURL); err != nil {
		return fmt.Errorf("invalid webhook URL: %w", err)
	}

	// Parse mentions from the message
	processedMessage, mobiles, isAtAll := parseMentions(message)

	// Check if message is in markdown format (starts with # or contains ** or __)
	isMarkdown := false
	if len(processedMessage) > 0 && (processedMessage[0] == '#' ||
		contains(processedMessage, "**") ||
		contains(processedMessage, "__") ||
		contains(processedMessage, "```") ||
		contains(processedMessage, ">")) {
		isMarkdown = true
	}

	// Create message payload
	msg := DingTalkMessage{}

	if isMarkdown {
		msg.MsgType = "markdown"
		msg.Markdown = map[string]string{
			"title": "Notification",
			"text":  processedMessage,
		}
	} else {
		msg.MsgType = "text"
		msg.Text = map[string]string{
			"content": processedMessage,
		}
	}

	// Add at configuration
	msg.At = map[string]interface{}{
		"atMobiles": mobiles,
		"isAtAll":   isAtAll,
	}

	// Convert message to JSON
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal DingTalk message: %w", err)
	}

	// Create request URL with signature if secret exists
	if d.Bot.Secret != "" {
		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		sign := calculateDingTalkSignature(timestamp, d.Bot.Secret)

		parsedURL, err := url.Parse(webhookURL)
		if err != nil {
			return fmt.Errorf("invalid webhook URL: %w", err)
		}

		q := parsedURL.Query()
		q.Set("timestamp", timestamp)
		q.Set("sign", sign)
		parsedURL.RawQuery = q.Encode()

		webhookURL = parsedURL.String()
	}

	// Send HTTP request
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		return fmt.Errorf("failed to send DingTalk message: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var response DingTalkResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode DingTalk response: %w", err)
	}

	// Check response status
	if response.ErrCode != 0 {
		return fmt.Errorf("failed to send DingTalk message: %s (code: %d)", response.ErrMsg, response.ErrCode)
	}

	return nil
}

// Helper function to check if string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// calculateDingTalkSignature calculates the signature for DingTalk webhook
func calculateDingTalkSignature(timestamp, secret string) string {
	stringToSign := timestamp + "\n" + secret
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
