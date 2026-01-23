package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type OpenRouterGenerator struct {
	client *http.Client
	apiKey string
	url    string
}

func NewOpenRouterGenerator(apiKey string) *OpenRouterGenerator {
	return &OpenRouterGenerator{
		client: &http.Client{Timeout: 30 * time.Second},
		apiKey: apiKey,
		url:    "",
	}
}

func (g *OpenRouterGenerator) withUrl(url string) *OpenRouterGenerator {
	newG := *g
	newG.url = url
	return &newG
}

func (g *OpenRouterGenerator) WithEmbeddings() *OpenRouterGenerator {
	return g.withUrl("https://api.openrouter.ai/v1/embeddings")
}

func (g *OpenRouterGenerator) WithCompletion() *OpenRouterGenerator {
	return g.withUrl("https://api.openrouter.ai/v1/chat/completions")
}

func (g *OpenRouterGenerator) Do(ctx context.Context, payload any, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if g.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+g.apiKey)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
