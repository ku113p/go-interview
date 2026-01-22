package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-interview/internal/history/domain"
)

var _ domain.UserHistoryList = (*Repository)(nil)
var _ domain.UserHistoryUpdater = (*Repository)(nil)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListUserHistory(ctx context.Context, userID uuid.UUID) ([]*domain.History, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, created_at, user_id, message
		FROM history
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list user history: %w", err)
	}
	defer rows.Close()

	items := make([]*domain.History, 0)
	for rows.Next() {
		var dto HistorySQL
		if err := dto.Scan(rows); err != nil {
			return nil, fmt.Errorf("scan history: %w", err)
		}

		historyItem, err := dto.ToDomain()
		if err != nil {
			return nil, err
		}
		items = append(items, historyItem)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate history rows: %w", err)
	}

	return items, nil
}

func (r *Repository) AddMessageToHistory(ctx context.Context, userID uuid.UUID, message domain.HistoryMessage, limit uint8) error {
	history := domain.NewHistory(uuid.New(), userID, message)
	dto, err := historyToSQL(history)
	if err != nil {
		return fmt.Errorf("history to dto: %w", err)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `
		INSERT INTO history (id, created_at, user_id, message)
		VALUES ($1, $2, $3, $4)
	`, dto.ID, dto.CreatedAt, dto.UserID, dto.Message); err != nil {
		return fmt.Errorf("insert history: %w", err)
	}

	if limit > 0 {
		if _, err = tx.Exec(ctx, `
			DELETE FROM history
			WHERE user_id = $1
			AND id IN (
				SELECT id FROM history
				WHERE user_id = $1
				ORDER BY created_at DESC
				OFFSET $2
			)
		`, userID, limit); err != nil {
			return fmt.Errorf("trim history: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit history tx: %w", err)
	}

	return nil
}
