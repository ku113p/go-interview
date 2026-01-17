package domain

import (
	"context"

	"github.com/google/uuid"
)

type LifeAreaRepository interface {
	Save(ctx context.Context, area *LifeArea) error
	GetByID(ctx context.Context, id uuid.UUID) (*LifeArea, error)
}
