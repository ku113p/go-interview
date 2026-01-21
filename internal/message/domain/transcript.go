package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transcript struct {
	ID uuid.UUID

	CreatedAt time.Time

	NodeID    uuid.UUID
	RawDataID uuid.UUID

	Content Content
}

func NewTranscript(id uuid.UUID, rawDataID uuid.UUID, nodeID uuid.UUID, content Content) *Transcript {
	return &Transcript{
		ID:        id,
		CreatedAt: time.Now(),
		NodeID:    nodeID,
		RawDataID: rawDataID,
		Content:   content,
	}
}
