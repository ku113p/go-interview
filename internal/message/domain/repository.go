package domain

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type RawDataCreator interface {
	CreateRawData(ctx context.Context, rd *RawData) error
}

type RawDataGetter interface {
	GetRawData(ctx context.Context, id uuid.UUID) (*RawData, error)
}

type RawDataDeleter interface {
	DeleteRawData(ctx context.Context, id uuid.UUID) error
}

type TranscriptCreator interface {
	CreateTranscript(ctx context.Context, t *Transcript) error
}

type TranscriptGetter interface {
	GetTranscript(ctx context.Context, id uuid.UUID) (*Transcript, error)
}

type TranscriptLister interface {
	ListTranscript(ctx context.Context, userID uuid.UUID) ([]*Transcript, error)
}

type TranscriptDeleter interface {
	DeleteTranscript(ctx context.Context, id uuid.UUID) error
}

type FileSaver interface {
	Save(ctx context.Context, stream io.Reader) (*string, error)
}

type FileGetter interface {
	Get(ctx context.Context, path string) (io.Reader, error)
}

type FileDeleter interface {
	Delete(ctx context.Context, path string) error
}

type AudioExtractor interface {
	ExtractAudio(ctx context.Context, stream io.Reader) (io.Reader, error)
}

type AudioToTextConverter interface {
	Convert(ctx context.Context, stream io.Reader) (string, error)
}

type IDGenerator interface {
	Generate() uuid.UUID
}
