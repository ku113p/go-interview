package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go-interview/internal/biography/domain"
	common "go-interview/internal/common/domain"
)

type Config struct {
	APIKey  string
	Model   string
	BaseURL string
	Timeout time.Duration
}

type OpenRouterGenerator struct {
	cfg    Config
	client *http.Client
}

func NewOpenRouterGenerator(cfg Config) *OpenRouterGenerator {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://openrouter.ai/api/v1/chat/completions"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.Model == "" {
		cfg.Model = "openai/gpt-4o-mini"
	}

	return &OpenRouterGenerator{
		cfg: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (g *OpenRouterGenerator) Generate(
	ctx context.Context,
	area *domain.LifeArea,
) ([]common.Description, error) {

	systemPrompt := "You are a helpful assistant that generates checklist criteria. You output ONLY raw JSON arrays of strings. No markdown formatting."
	userPrompt := fmt.Sprintf(
		`Generate a checklist of specific, actionable criteria for the following life area.
Title: %s
Goal: %s

Output format: ["criteria 1", "criteria 2", "criteria 3"]`,
		area.Title,
		area.Goal,
	)

	reqBody := openRouterRequest{
		Model: g.cfg.Model,
		Messages: []message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.cfg.BaseURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+g.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "http://localhost")
	req.Header.Set("X-Title", "biography-service")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("provider returned status %d: %s", resp.StatusCode, string(errorBody))
	}

	var orResp openRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&orResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(orResp.Choices) == 0 {
		return nil, fmt.Errorf("provider returned no choices")
	}

	rawContent := orResp.Choices[0].Message.Content
	cleanedContent := cleanJSONMarkdown(rawContent)

	var items []string
	if err := json.Unmarshal([]byte(cleanedContent), &items); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON: %w (content: %s)", err, rawContent)
	}

	result := make([]common.Description, 0, len(items))
	for _, s := range items {
		result = append(result, common.NewDescription(s))
	}

	return result, nil
}

func cleanJSONMarkdown(input string) string {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "```json") {
		input = strings.TrimPrefix(input, "```json")
		input = strings.TrimSuffix(input, "```")
	} else if strings.HasPrefix(input, "```") {
		input = strings.TrimPrefix(input, "```")
		input = strings.TrimSuffix(input, "```")
	}
	return strings.TrimSpace(input)
}

type openRouterRequest struct {
	Model          string          `json:"model"`
	Messages       []message       `json:"messages"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}
