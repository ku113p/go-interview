package domain

import (
	"context"

	"github.com/google/uuid"
)

type FactCreator interface {
	CreateFacts(ctx context.Context, facts []*Fact) error
}

type EmbeddingCreator interface {
	CreateEmbeddings(ctx context.Context, embeddings []*Embedding) error
}

type VectorGenerator interface {
	GenerateVector(ctx context.Context, text string) (Vector, error)
}

type IDGenerator interface {
	Generate() (uuid.UUID, error)
}

type FactExtractor interface {
	ExtractFacts(ctx context.Context, text, goal string) ([]string, error)
}
