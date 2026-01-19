package postgres

import (
	"time"

	"github.com/google/uuid"
)

type LifeAreaSQL struct {
	ID        uuid.UUID       `db:"id"`
	ParentID  uuid.UUID       `db:"parent_id"`
	UserID    uuid.UUID       `db:"user_id"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
	Title     string          `db:"title"`
	Goal      string          `db:"goal"`
	Criteria  []*CriterionSQL `db:"criteria"`
}

type CriterionSQL struct {
	ID          uuid.UUID `db:"id"`
	NodeID      uuid.UUID `db:"node_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Description string    `db:"description"`
	IsCompleted bool      `db:"is_completed"`
}
