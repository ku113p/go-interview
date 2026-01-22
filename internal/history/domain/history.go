package domain

import (
	"time"

	"github.com/google/uuid"
)

type HistoryMessage map[string]interface{}

type History struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UserID    uuid.UUID
	Message   map[string]interface{}
}

func NewHistory(id uuid.UUID, userID uuid.UUID, message HistoryMessage) *History {
	return &History{
		ID:        id,
		CreatedAt: time.Now(),
		UserID:    userID,
		Message:   message,
	}
}
