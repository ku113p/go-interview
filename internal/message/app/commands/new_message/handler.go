package newmessage

import (
	"context"
	"fmt"
	"go-interview/internal/message/domain"
	"io"

	"github.com/google/uuid"
)

type Repository interface {
	domain.RawDataCreator
	domain.TranscriptCreator
}

type FileStorage interface {
	domain.FileSaver
}

type MediaProcessor interface {
	domain.AudioExtractor
}

type AiClient interface {
	domain.AudioToTextConverter
}

type TranscriptionHandler struct {
	repo      Repository
	storage   FileStorage
	processor MediaProcessor
	aiClient  AiClient
	genID     domain.IDGenerator
}

func NewTranscriptionHandler(r Repository, s FileStorage, p MediaProcessor, ai AiClient, genID domain.IDGenerator) *TranscriptionHandler {
	return &TranscriptionHandler{
		repo:      r,
		storage:   s,
		processor: p,
		aiClient:  ai,
		genID:     genID,
	}
}

func (h *TranscriptionHandler) Handle(ctx context.Context, cmd *NewMessageCommand) (*NewMessageResult, error) {
	nodeID, err := uuid.Parse(cmd.NodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse node id: %w", err)
	}

	var extractedText string
	if cmd.MediaType == "text" {
		b, err := io.ReadAll(cmd.Stream)
		if err != nil {
			return nil, fmt.Errorf("failed to read stream: %w", err)
		}
		extractedText = string(b)
	} else {
		var audioStream io.Reader
		var err error
		if cmd.MediaType == "audio" {
			audioStream, err = h.processor.ExtractAudio(ctx, cmd.Stream)
			if err != nil {
				return nil, fmt.Errorf("failed to extract audio: %w", err)
			}
		} else if cmd.MediaType == "video" {
			audioStream, err = h.processor.ExtractAudio(ctx, cmd.Stream)
			if err != nil {
				return nil, fmt.Errorf("failed to extract audio: %w", err)
			}
		} else {
			return nil, fmt.Errorf("unsupported media type: %s", cmd.MediaType)
		}

		text, err := h.aiClient.Convert(ctx, audioStream)
		if err != nil {
			return nil, fmt.Errorf("failed to convert audio to text: %w", err)
		}
		extractedText = text
	}

	rawDataStoragePath, err := h.storage.Save(ctx, cmd.Stream)
	if err != nil {
		return nil, err
	}

	rawData := domain.NewRawData(
		h.genID.Generate(),
		rawDataStoragePath,
		domain.MediaType(cmd.MediaType),
	)
	err = h.repo.CreateRawData(ctx, rawData)
	if err != nil {
		return nil, err
	}

	transcript := domain.NewTranscript(
		h.genID.Generate(),
		rawData.ID,
		nodeID,
		domain.NewContent(extractedText),
	)
	err = h.repo.CreateTranscript(ctx, transcript)
	if err != nil {
		return nil, err
	}

	return &NewMessageResult{
		Text: extractedText,
	}, nil
}
