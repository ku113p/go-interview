package queries

import (
	"context"
	"go-interview/internal/history/domain"

	"github.com/google/uuid"
)

type Repository interface {
	domain.UserHistoryList
}

type GetHistoryHandler struct {
	repo Repository
}

func NewGetHistoryHandler(repo Repository) *GetHistoryHandler {
	return &GetHistoryHandler{
		repo: repo,
	}
}

func (h *GetHistoryHandler) Handle(ctx context.Context, query GetHistoryQuery) (*GetHistoryResult, error) {
	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, err
	}

	items, err := h.repo.ListUserHistory(ctx, userID)
	if err != nil {
		return nil, err
	}

	resultItems := make([]*HistoryMessage, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		msg := HistoryMessage(item.Message)
		resultItems = append(resultItems, &msg)
	}

	return &GetHistoryResult{
		Items: resultItems,
	}, nil
}
