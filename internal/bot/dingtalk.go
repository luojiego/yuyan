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

	// Check if message is in markdown format (starts with # or contains ** or __)
	isMarkdown := false
	if len(message) > 0 && (message[0] == '#' ||
		contains(message, "**") ||
		contains(message, "__") ||
		contains(message, "```") ||
		contains(message, ">")) {
		isMarkdown = true
	}

	// Create message payload
	msg := DingTalkMessage{}

	if isMarkdown {
		msg.MsgType = "markdown"
		msg.Markdown = map[string]string{
			"title": "Notification",
			"text":  message,
		}
	} else {
		msg.MsgType = "text"
		msg.Text = map[string]string{
			"content": message,
		}
	}

	// Add at configuration if needed
	// This could be enhanced to parse @mentions from the message
	// msg.At = map[string]interface{}{
	//    "atMobiles": []string{},
	//    "isAtAll": false,
	// }

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
