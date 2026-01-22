package postgres

import (
	"encoding/json"
	"fmt"

	"go-interview/internal/history/domain"
)

func historyToSQL(history *domain.History) (*HistorySQL, error) {
	message, err := json.Marshal(history.Message)
	if err != nil {
		return nil, fmt.Errorf("marshal history message: %w", err)
	}

	return &HistorySQL{
		ID:        history.ID,
		CreatedAt: history.CreatedAt,
		UserID:    history.UserID,
		Message:   message,
	}, nil
}

func (dto *HistorySQL) ToDomain() (*domain.History, error) {
	var message domain.HistoryMessage
	if err := json.Unmarshal(dto.Message, &message); err != nil {
		return nil, fmt.Errorf("unmarshal history message: %w", err)
	}

	return &domain.History{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		UserID:    dto.UserID,
		Message:   message,
	}, nil
}
