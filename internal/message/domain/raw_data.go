package domain

import (
	"time"

	"github.com/google/uuid"
)

type MediaType string

const (
	Video MediaType = "video"
	Audio MediaType = "audio"
	Text  MediaType = "text"
)

type RawData struct {
	ID uuid.UUID

	CreatedAt time.Time

	S3Path    string
	MediaType MediaType
}

func NewRawData(id uuid.UUID, s3Path string, mediaType MediaType) *RawData {
	return &RawData{
		ID:        id,
		CreatedAt: time.Now(),
		S3Path:    s3Path,
		MediaType: mediaType,
	}
}
