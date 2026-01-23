package aiservice

import (
	"context"
	"encoding/json"
	"fmt"

	common "go-interview/internal/common/domain"
	"go-interview/internal/criterion/domain"
)

type OpenRouterCriteriaApprover struct {
	Generator *common.OpenRouterGenerator
}

func NewOpenRouterCriteriaApprover(generator *common.OpenRouterGenerator) *OpenRouterCriteriaApprover {
	return &OpenRouterCriteriaApprover{
		Generator: generator,
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type chatRequest struct {
	Model          string         `json:"model"`
	Messages       []chatMessage  `json:"messages"`
	ResponseFormat responseFormat `json:"response_format"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type meetsResponse struct {
	Results []meetResult `json:"results"`
}

type meetResult struct {
	Description string `json:"description"`
	Meets       bool   `json:"meets"`
}

func (s *OpenRouterCriteriaApprover) MeetCriteria(ctx context.Context, text string, criteria []*domain.Criterion) ([]bool, error) {
	systemPrompt := `You are a strict compliance checker. 
	Output ONLY valid JSON.
	Schema: {"results": [{"description": "exact criteria text", "meets": true/false}]}`

	criteriaJSON, _ := json.Marshal(descriptions)
	userContent := fmt.Sprintf("Text:\n%s\n\nCheck against these criteria:\n%s", text, string(criteriaJSON))

	reqPayload := chatRequest{
		Model: "openai/gpt-4o-mini",
		ResponseFormat: responseFormat{
			Type: "json_object",
		},
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userContent},
		},
	}

	var apiResp chatResponse

	err := s.Generator.WithCompletion().Do(ctx, reqPayload, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("generator request failed: %w", err)
	}

	if apiResp.Error != nil {
		return nil, fmt.Errorf("provider api error: %s", apiResp.Error.Message)
	}
	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("provider returned no choices")
	}

	var aiResult meetsResponse
	contentString := apiResp.Choices[0].Message.Content

	if err := json.Unmarshal([]byte(contentString), &aiResult); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON content: %w", err)
	}

	resultsMap := make(map[string]bool)
	for _, res := range aiResult.Results {
		resultsMap[res.Description] = res.Meets
	}

	finalOutput := make([]bool, len(descriptions))
	for i, desc := range descriptions {
		if val, exists := resultsMap[desc]; exists {
			finalOutput[i] = val
		} else {
			finalOutput[i] = false
		}
	}

	return finalOutput, nil
}
