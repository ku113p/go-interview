package domain

import (
	"context"

	"github.com/google/uuid"
)

type LifeAreaRepository interface {
	CreateLifeArea(ctx context.Context, la *LifeArea) error
	GetLifeArea(ctx context.Context, ID uuid.UUID) (*LifeArea, error)
	ListLifeAreas(ctx context.Context, userID uuid.UUID) ([]*LifeArea, error)
	DeleteLifeArea(ctx context.Context, ID uuid.UUID) error
	UpdateLifeArea(ctx context.Context, la *LifeArea) error
}
