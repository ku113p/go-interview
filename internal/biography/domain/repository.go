package domain

import (
	"context"

	"github.com/google/uuid"
)

type LifeAreaCreator interface {
	CreateLifeArea(ctx context.Context, la *LifeArea) error
}

type LifeAreaGetter interface {
	GetLifeArea(ctx context.Context, ID uuid.UUID) (*LifeArea, error)
}

type LifeAreaLister interface {
	ListLifeAreas(ctx context.Context, userID uuid.UUID) ([]*LifeArea, error)
}

type LifeAreaDeleter interface {
	DeleteLifeArea(ctx context.Context, ID uuid.UUID) error
}

type LifeAreaParentChanger interface {
	ChangeParentID(ctx context.Context, ID uuid.UUID, parentID *uuid.UUID) error
}

type LifeAreaGoalChanger interface {
	ChangeGoal(ctx context.Context, ID uuid.UUID, goal Goal) error
}

type CriteriaDeleter interface {
	DeleteCriteria(ctx context.Context, IDs ...uuid.UUID) error
}

type CriteriaCreator interface {
	CreateCriteria(ctx context.Context, criteria ...*Criterion) error
}

type CriteriaNodeGetter interface {
	GetCriteriaNodeIDs(ctx context.Context, IDs ...uuid.UUID) (map[uuid.UUID]uuid.UUID, error)
}

type IDGenerator interface {
	GenID() (uuid.UUID, error)
}
