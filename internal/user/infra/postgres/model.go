package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserSQL struct {
	ID         uuid.UUID `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	ExternalID string    `db:"external_id"`
}

func (dto *UserSQL) Scan(row pgx.Row) error {
	return row.Scan(
		&dto.ID,
		&dto.CreatedAt,
		&dto.ExternalID,
	)
}
