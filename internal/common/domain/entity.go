package domain

import (
	"time"

	"github.com/google/uuid"
)

type Entity struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

func InitEntity(entity *Entity) {
	entity.ID = uuid.New()
	entity.CreatedAt = time.Now()
}

type UpdatableEntity struct {
	Entity
	UpdatedAt time.Time
}

func (e *UpdatableEntity) Update() {
	e.UpdatedAt = time.Now()
}

func InitUpdatableEntity(entity *UpdatableEntity) {
	now := time.Now()
	entity.ID = uuid.New()
	entity.CreatedAt = now
	entity.UpdatedAt = now
}
