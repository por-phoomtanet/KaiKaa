package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Completer — interface สำหรับเรียก LLM (แยกไว้เพื่อ mock ใน test ได้ ไม่ต้องใช้ key จริง)
type Completer interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

// OpenRouterClient เรียก openrouter chat completions
type OpenRouterClient struct {
	apiKey string
	model  string
	url    string
	http   *http.Client
}

func NewOpenRouterClient(apiKey, model string) *OpenRouterClient {
	return &OpenRouterClient{
		apiKey: apiKey,
		model:  model,
		url:    "https://openrouter.ai/api/v1/chat/completions",
		http:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *OpenRouterClient) Complete(ctx context.Context, prompt string) (string, error) {
	if c.apiKey == "" {
		return "", errors.New("OPENROUTER_KEY not configured")
	}

	payload, _ := json.Marshal(map[string]any{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("openrouter request failed: %w", err)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openrouter status %d: %s", resp.StatusCode, string(data))
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", fmt.Errorf("openrouter decode failed: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return "", errors.New("openrouter returned no choices")
	}
	return parsed.Choices[0].Message.Content, nil
}
