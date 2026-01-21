package createembeddings

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go-interview/internal/memory/domain"
)

const (
	defaultChunkSize = 800
	chunkOverlap     = 100
)

type Repository interface {
	domain.EmbeddingCreator
}

type VectorGenerator interface {
	domain.VectorGenerator
}

type Handler struct {
	repo   Repository
	genID  domain.IDGenerator
	genVec VectorGenerator
}

func NewHandler(repo Repository, genID domain.IDGenerator, genVec VectorGenerator) *Handler {
	return &Handler{
		repo:   repo,
		genID:  genID,
		genVec: genVec,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	text := strings.TrimSpace(cmd.Text)
	if text == "" {
		return errors.New("text is required")
	}

	chunks := splitText(text, defaultChunkSize, chunkOverlap)
	if len(chunks) == 0 {
		return errors.New("text splitting resulted in no chunks")
	}

	embeddings := make([]*domain.Embedding, 0, len(chunks))
	for _, chunk := range chunks {
		vector, err := h.genVec.GenerateVector(ctx, chunk)
		if err != nil {
			return fmt.Errorf("generate vector: %w", err)
		}

		embID, err := h.genID.Generate()
		if err != nil {
			return fmt.Errorf("generate id: %w", err)
		}

		embeddings = append(embeddings, domain.NewEmbedding(
			embID,
			cmd.NodeID,
			vector,
			domain.NewContent(chunk),
		))
	}

	if err := h.repo.CreateEmbeddings(ctx, embeddings); err != nil {
		return fmt.Errorf("create embeddings: %w", err)
	}

	return nil
}

func splitText(text string, chunkSize, overlap int) []string {
	if chunkSize <= 0 {
		return nil
	}

	var chunks []string
	runes := []rune(text)
	length := len(runes)
	start := 0

	for start < length {
		end := start + chunkSize
		if end > length {
			end = length
		}

		chunk := strings.TrimSpace(string(runes[start:end]))
		if chunk != "" {
			chunks = append(chunks, chunk)
		}

		if end == length {
			break
		}

		start = end - overlap
		if start < 0 {
			start = 0
		}
	}

	return chunks
}
