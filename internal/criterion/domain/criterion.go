package domain

import (
	"time"

	"github.com/google/uuid"
)

type Criterion struct {
	ID uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time

	NodeID uuid.UUID

	Description Description
	IsCompleted bool
}

func NewCriterion(id uuid.UUID, nodeID uuid.UUID, description Description) *Criterion {
	now := time.Now()

	return &Criterion{
		ID:          id,
		CreatedAt:   now,
		UpdatedAt:   now,
		NodeID:      nodeID,
		Description: description,
		IsCompleted: false,
	}
}
