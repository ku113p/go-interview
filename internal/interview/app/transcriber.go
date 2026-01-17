package app

import (
	"context"
	"fmt"
	common "go-interview/internal/common/domain"
	"go-interview/internal/interview/domain"
	"io"

	"github.com/google/uuid"
)

type Repository interface {
	GetRawData(ctx context.Context, id uuid.UUID) (*domain.RawData, error)
	SaveTranscript(ctx context.Context, t *domain.Transcript) error
}

type FileStorage interface {
	Download(ctx context.Context, path string) (io.Reader, error)
}

type MediaProcessor interface {
	ExtractAudio(ctx context.Context, mediaStream io.Reader) (io.Reader, error)
}

type AiClient interface {
	SpeechToText(ctx context.Context, audioData io.Reader) (string, error)
}

type TranscriptionService struct {
	repo      Repository
	storage   FileStorage
	processor MediaProcessor
	aiClient  AiClient
}

func NewTranscriptionService(r Repository, s FileStorage, p MediaProcessor, ai AiClient) *TranscriptionService {
	return &TranscriptionService{
		repo:      r,
		storage:   s,
		processor: p,
		aiClient:  ai,
	}
}

func (s *TranscriptionService) Transcribe(ctx context.Context, rawDataID uuid.UUID) error {
	rawData, err := s.repo.GetRawData(ctx, rawDataID)
	if err != nil {
		return fmt.Errorf("failed to get raw data: %w", err)
	}

	fileStream, err := s.storage.Download(ctx, rawData.S3Path)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	var recognizedText string
	switch rawData.MediaType {
	case domain.Text:
		b, err := io.ReadAll(fileStream)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		recognizedText = string(b)
	default:
		var audioStream io.Reader
		switch rawData.MediaType {
		case domain.Video:
			audioStream, err = s.processor.ExtractAudio(ctx, fileStream)
			if err != nil {
				return fmt.Errorf("failed to extract audio from video: %w", err)
			}
		case domain.Audio:
			audioStream = fileStream
		default:
			return fmt.Errorf("unsupported media type: %s", rawData.MediaType)
		}
		recognizedText, err = s.aiClient.SpeechToText(ctx, audioStream)
		if err != nil {
			return fmt.Errorf("failed to transcribe audio: %w", err)
		}
	}

	content := common.NewContent(recognizedText)
	transcript := domain.NewTranscript(rawDataID, content)

	if err := s.repo.SaveTranscript(ctx, transcript); err != nil {
		return fmt.Errorf("failed to save transcript: %w", err)
	}

	return nil
}
