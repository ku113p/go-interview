package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RawDataSQL struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	S3Path    string    `db:"s3_path"`
	MediaType string    `db:"media_type"`
}

func (dto *RawDataSQL) Scan(row pgx.Row) error {
	return row.Scan(
		&dto.ID,
		&dto.CreatedAt,
		&dto.S3Path,
		&dto.MediaType,
	)
}

type TranscriptSQL struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	NodeID    uuid.UUID `db:"node_id"`
	RawDataID uuid.UUID `db:"raw_data_id"`
	Content   string    `db:"content"`
}

func (dto *TranscriptSQL) Scan(row pgx.Row) error {
	return row.Scan(
		&dto.ID,
		&dto.CreatedAt,
		&dto.NodeID,
		&dto.RawDataID,
		&dto.Content,
	)
}
