package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-interview/internal/memory/domain"
)

const (
	defaultFactsModel = "anthropic/claude-3.5-sonnet"
	defaultFactsURL   = "https://openrouter.ai/api/v1/chat/completions"
)

type OpenRouterFactExtractor struct {
	client *http.Client
	apiKey string
	model  string
	url    string
}

func NewOpenRouterFactExtractor(apiKey, model, url string) *OpenRouterFactExtractor {
	if model == "" {
		model = defaultFactsModel
	}
	if url == "" {
		url = defaultFactsURL
	}

	return &OpenRouterFactExtractor{
		client: &http.Client{Timeout: 30 * time.Second},
		apiKey: apiKey,
		model:  model,
		url:    url,
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Response string        `json:"response_format,omitempty"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type factsResponse struct {
	Facts []string `json:"facts"`
}

func (e *OpenRouterFactExtractor) ExtractFacts(ctx context.Context, text, goal string) ([]string, error) {
	prompt := buildFactsPrompt(text, goal)
	payload := chatRequest{
		Model: e.model,
		Messages: []chatMessage{
			{Role: "system", Content: "You are a careful analyst that extracts factual statements from text."},
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if e.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+e.apiKey)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var out chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if out.Error != nil {
		return nil, fmt.Errorf("api error: %s", out.Error.Message)
	}

	if len(out.Choices) == 0 {
		return nil, errors.New("empty choices in response")
	}

	content := strings.TrimSpace(out.Choices[0].Message.Content)
	if content == "" {
		return nil, errors.New("empty content in response")
	}

	var parsed factsResponse
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, fmt.Errorf("parse facts json: %w", err)
	}

	var facts []string
	for _, fact := range parsed.Facts {
		trimmed := strings.TrimSpace(fact)
		if trimmed != "" {
			facts = append(facts, trimmed)
		}
	}

	return facts, nil
}

func buildFactsPrompt(text, goal string) string {
	var b strings.Builder
	b.WriteString("You will receive text and a goal. Extract concise, self-contained factual statements that directly support the goal. Return only JSON: {\"facts\":[\"fact1\",...]}\n")
	b.WriteString("Goal:\n")
	b.WriteString(goal)
	b.WriteString("\n\nText:\n")
	b.WriteString(text)
	return b.String()
}

var _ domain.FactExtractor = (*OpenRouterFactExtractor)(nil)
