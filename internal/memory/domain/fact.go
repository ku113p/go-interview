package domain

import (
	"time"

	"github.com/google/uuid"
)

type Fact struct {
	ID uuid.UUID

	CreatedAt time.Time

	NodeID uuid.UUID

	Info     Info
	DateTime *time.Time
}

func NewFact(id uuid.UUID, nodeID uuid.UUID, info Info, dateTime *time.Time) *Fact {
	return &Fact{
		ID:        id,
		CreatedAt: time.Now(),
		NodeID:    nodeID,
		Info:      info,
		DateTime:  dateTime,
	}
}
