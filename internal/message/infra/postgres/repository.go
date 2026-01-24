package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	common "go-interview/internal/common/domain"
	"go-interview/internal/message/domain"
)

var _ domain.RawDataCreator = (*PostgresRepository)(nil)
var _ domain.RawDataGetter = (*PostgresRepository)(nil)
var _ domain.RawDataDeleter = (*PostgresRepository)(nil)
var _ domain.TranscriptCreator = (*PostgresRepository)(nil)
var _ domain.TranscriptGetter = (*PostgresRepository)(nil)
var _ domain.TranscriptLister = (*PostgresRepository)(nil)
var _ domain.TranscriptDeleter = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateRawData(ctx context.Context, rd *domain.RawData) error {
	dto := rawDataToSQL(rd)

	_, err := r.db.Exec(ctx, `
		INSERT INTO raw_data (id, created_at, s3_path, media_type)
		VALUES ($1, $2, $3, $4)
	`, dto.ID, dto.CreatedAt, dto.S3Path, dto.MediaType)
	if err != nil {
		return fmt.Errorf("insert raw data: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetRawData(ctx context.Context, id uuid.UUID) (*domain.RawData, error) {
	var dto RawDataSQL
	row := r.db.QueryRow(ctx, `
		SELECT id, created_at, s3_path, media_type
		FROM raw_data
		WHERE id = $1
	`, id)
	if err := dto.Scan(row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common.ErrNotFound
		}
		return nil, fmt.Errorf("get raw data: %w", err)
	}

	return dto.ToDomain(), nil
}

func (r *PostgresRepository) DeleteRawData(ctx context.Context, id uuid.UUID) error {
	if _, err := r.db.Exec(ctx, `
		DELETE FROM raw_data
		WHERE id = $1
	`, id); err != nil {
		return fmt.Errorf("delete raw data: %w", err)
	}

	return nil
}

func (r *PostgresRepository) CreateTranscript(ctx context.Context, t *domain.Transcript) error {
	dto := transcriptToSQL(t)

	_, err := r.db.Exec(ctx, `
		INSERT INTO transcripts (id, created_at, node_id, raw_data_id, content)
		VALUES ($1, $2, $3, $4, $5)
	`, dto.ID, dto.CreatedAt, dto.NodeID, dto.RawDataID, dto.Content)
	if err != nil {
		return fmt.Errorf("insert transcript: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetTranscript(ctx context.Context, id uuid.UUID) (*domain.Transcript, error) {
	var dto TranscriptSQL
	row := r.db.QueryRow(ctx, `
		SELECT id, created_at, node_id, raw_data_id, content
		FROM transcripts
		WHERE id = $1
	`, id)
	if err := dto.Scan(row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common.ErrNotFound
		}
		return nil, fmt.Errorf("get transcript: %w", err)
	}

	return dto.ToDomain(), nil
}

func (r *PostgresRepository) ListTranscript(ctx context.Context, userID uuid.UUID) ([]*domain.Transcript, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, created_at, node_id, raw_data_id, content
		FROM transcripts
		WHERE node_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list transcripts: %w", err)
	}
	defer rows.Close()

	items := make([]*domain.Transcript, 0)
	for rows.Next() {
		var dto TranscriptSQL
		if err := dto.Scan(rows); err != nil {
			return nil, fmt.Errorf("scan transcript: %w", err)
		}
		items = append(items, dto.ToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate transcripts: %w", err)
	}

	return items, nil
}

func (r *PostgresRepository) DeleteTranscript(ctx context.Context, id uuid.UUID) error {
	if _, err := r.db.Exec(ctx, `
		DELETE FROM transcripts
		WHERE id = $1
	`, id); err != nil {
		return fmt.Errorf("delete transcript: %w", err)
	}

	return nil
}
