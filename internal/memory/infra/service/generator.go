package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-interview/internal/memory/domain"
)

const defaultEmbeddingsURL = "https://openrouter.ai/api/v1/embeddings"

type OpenRouterGenerator struct {
	client *http.Client
	apiKey string
	model  string
	url    string
}

func NewOpenRouterGenerator(apiKey, model, url string) *OpenRouterGenerator {
	if model == "" {
		model = "text-embedding-3-small"
	}
	if url == "" {
		url = defaultEmbeddingsURL
	}

	return &OpenRouterGenerator{
		client: &http.Client{Timeout: 30 * time.Second},
		apiKey: apiKey,
		model:  model,
		url:    url,
	}
}

type embeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (g *OpenRouterGenerator) GenerateVector(ctx context.Context, text string) (domain.Vector, error) {
	payload := embeddingRequest{
		Model: g.model,
		Input: text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if g.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+g.apiKey)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var out embeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if out.Error != nil {
		return nil, fmt.Errorf("api error: %s", out.Error.Message)
	}

	if len(out.Data) == 0 {
		return nil, fmt.Errorf("empty embedding response")
	}

	return domain.NewVector(out.Data[0].Embedding), nil
}
