package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	common "go-interview/internal/common/domain"
	"go-interview/internal/memory/domain"
)

type OpenRouterFactExtractor struct {
	Generator common.OpenRouterGenerator
}

func NewOpenRouterFactExtractor(generator common.OpenRouterGenerator) *OpenRouterFactExtractor {
	return &OpenRouterFactExtractor{
		Generator: generator,
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
		Model: "anthropic/claude-3.5-sonnet",
		Messages: []chatMessage{
			{Role: "system", Content: "You are a careful analyst that extracts factual statements from text."},
			{Role: "user", Content: prompt},
		},
	}

	var out chatResponse
	if err := e.Generator.WithCompletion().Do(ctx, payload, &out); err != nil {
		return nil, err
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
