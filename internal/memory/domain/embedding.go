package domain

import (
	"time"

	"github.com/google/uuid"
)

type Embedding struct {
	ID uuid.UUID

	CreatedAt time.Time

	NodeID uuid.UUID

	Vector  Vector
	Content Content
}

func NewEmbedding(id uuid.UUID, nodeID uuid.UUID, vector Vector, content Content) *Embedding {
	return &Embedding{
		ID:        id,
		NodeID:    nodeID,
		CreatedAt: time.Now(),
		Vector:    vector,
		Content:   content,
	}
}
