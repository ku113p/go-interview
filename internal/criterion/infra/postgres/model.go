package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CriterionSQL struct {
	ID          uuid.UUID `db:"id"`
	NodeID      uuid.UUID `db:"node_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Description string    `db:"description"`
	IsCompleted bool      `db:"is_completed"`
}

func (dto *CriterionSQL) Scan(s pgx.Row) error {
	return s.Scan(
		&dto.ID,
		&dto.NodeID,
		&dto.CreatedAt,
		&dto.UpdatedAt,
		&dto.Description,
		&dto.IsCompleted,
	)
}
