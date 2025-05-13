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

// NewDingTalkBot creates a new DingTalk bot instance
func NewDingTalkBot(bot models.Bot) *DingTalkBot {
	return &DingTalkBot{
		Bot: bot,
	}
}

// Send sends a message through the DingTalk bot
func (d *DingTalkBot) Send(message string) error {
	// Create message payload
	msg := DingTalkMessage{
		MsgType: "text",
		Text: map[string]string{
			"content": message,
		},
	}

	// Convert message to JSON
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal DingTalk message: %w", err)
	}

	// Create request URL with signature if secret exists
	webhookURL := d.Bot.WebhookURL
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

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("failed to send DingTalk message: status %d", resp.StatusCode)
		}
		return fmt.Errorf("failed to send DingTalk message: %v", result)
	}

	return nil
}

// calculateDingTalkSignature calculates the signature for DingTalk webhook
func calculateDingTalkSignature(timestamp, secret string) string {
	stringToSign := timestamp + "\n" + secret
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
