package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserHistoryList interface {
	ListUserHistory(ctx context.Context, userID uuid.UUID) ([]*History, error)
}

type UserHistoryUpdater interface {
	AddMessageToHistory(ctx context.Context, userID uuid.UUID, message HistoryMessage, limit uint8) error
}
