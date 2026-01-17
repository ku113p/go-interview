package domain

import (
	common "go-interview/internal/common/domain"

	"github.com/google/uuid"
)

type Transcript struct {
	common.UpdatableEntity
	RawDataID uuid.UUID
	Content   common.Content
}

func NewTranscript(rawDataID uuid.UUID, content common.Content) *Transcript {
	transcript := &Transcript{
		RawDataID: rawDataID,
		Content:   content,
	}

	common.InitUpdatableEntity(&transcript.UpdatableEntity)

	return transcript
}
