package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FactSQL struct {
	ID        uuid.UUID  `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	NodeID    uuid.UUID  `db:"node_id"`
	Info      string     `db:"info"`
	DateTime  *time.Time `db:"date_time"`
}

type EmbeddingSQL struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	NodeID    uuid.UUID `db:"node_id"`
	Vector    []float64 `db:"vector"`
	Content   string    `db:"content"`
}

func (dto *FactSQL) Scan(row pgx.Row) error {
	return row.Scan(
		&dto.ID,
		&dto.CreatedAt,
		&dto.NodeID,
		&dto.Info,
		&dto.DateTime,
	)
}

func (dto *EmbeddingSQL) Scan(row pgx.Row) error {
	return row.Scan(
		&dto.ID,
		&dto.CreatedAt,
		&dto.NodeID,
		&dto.Vector,
		&dto.Content,
	)
}
