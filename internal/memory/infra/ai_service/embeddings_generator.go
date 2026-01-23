package service

import (
	"context"
	"fmt"

	common "go-interview/internal/common/domain"
	"go-interview/internal/memory/domain"
)

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

type OpenRouterEmbGenerator struct {
	Generator common.OpenRouterGenerator
}

func NewAIService(generator common.OpenRouterGenerator) *OpenRouterEmbGenerator {
	return &OpenRouterEmbGenerator{
		Generator: generator,
	}
}

func (s *OpenRouterEmbGenerator) GenerateVector(ctx context.Context, text string) (domain.Vector, error) {
	payload := embeddingRequest{
		Model: "text-embedding-3-small",
		Input: text,
	}

	var out embeddingResponse
	if err := s.Generator.WithEmbeddings().Do(ctx, payload, &out); err != nil {
		return nil, err
	}

	if out.Error != nil {
		return nil, fmt.Errorf("api error: %s", out.Error.Message)
	}

	if len(out.Data) == 0 {
		return nil, fmt.Errorf("empty embedding response")
	}

	return domain.NewVector(out.Data[0].Embedding), nil
}
