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
	return &Criterion{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		NodeID:      nodeID,
		Description: description,
		IsCompleted: false,
	}
}
