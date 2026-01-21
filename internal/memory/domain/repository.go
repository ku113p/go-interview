package domain

import (
	"context"

	"github.com/google/uuid"
)

type FactCreator interface {
	CreateFact(ctx context.Context, fact *Fact) error
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
