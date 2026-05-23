package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.telegram.org"

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
}

type SendParams struct {
	Token  string
	ChatID string
	Text   string
}

type apiResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		BaseURL:    defaultBaseURL,
	}
}

func (c *Client) Send(p SendParams) error {
	if p.Token == "" {
		return fmt.Errorf("bot token is required")
	}
	if p.ChatID == "" {
		return fmt.Errorf("chat ID is required")
	}
	if p.Text == "" {
		return fmt.Errorf("message text is required")
	}

	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	apiURL := baseURL + "/bot" + p.Token + "/sendMessage"

	data := url.Values{}
	data.Set("chat_id", p.ChatID)
	data.Set("text", p.Text)

	resp, err := c.HTTPClient.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !apiResp.OK {
		msg := apiResp.Description
		if msg == "" {
			msg = "unknown error"
		}
		return fmt.Errorf("telegram API error: %s", msg)
	}

	return nil
}
