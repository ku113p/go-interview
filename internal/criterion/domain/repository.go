package domain

import (
	"context"

	"github.com/google/uuid"
)

type CriteriaDeleter interface {
	DeleteCriteria(ctx context.Context, IDs ...uuid.UUID) error
}

type CriteriaCreator interface {
	CreateCriteria(ctx context.Context, criteria ...*Criterion) error
}

type CriteriaLister interface {
	GetCriteriaByNode(ctx context.Context, nodeID uuid.UUID) ([]*Criterion, error)
}

type CriteriaMarker interface {
	MarkCriteria(ctx context.Context, IDs ...uuid.UUID) error
}

type IDGenerator interface {
	Generate() (uuid.UUID, error)
}
