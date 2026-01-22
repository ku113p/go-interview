package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type HistorySQL struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UserID    uuid.UUID `db:"user_id"`
	Message   []byte    `db:"message"`
}

func (dto *HistorySQL) Scan(row pgx.Row) error {
	return row.Scan(
		&dto.ID,
		&dto.CreatedAt,
		&dto.UserID,
		&dto.Message,
	)
}
