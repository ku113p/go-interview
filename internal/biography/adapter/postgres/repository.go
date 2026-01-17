package postgres

import (
	"context"
	"fmt"
	"go-interview/internal/biography/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(ctx context.Context, area *domain.LifeArea) error {
	lifeAreaSQL, criterionSQL := toSQL(area)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO nodes (id, parent_id, user_id, created_at, updated_at, title, goal)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, lifeAreaSQL.ID, lifeAreaSQL.ParentID, lifeAreaSQL.UserID, lifeAreaSQL.CreatedAt, lifeAreaSQL.UpdatedAt, lifeAreaSQL.Title, lifeAreaSQL.Goal)
	if err != nil {
		return fmt.Errorf("save node: %w", err)
	}

	for _, c := range criterionSQL {
		_, err = tx.Exec(ctx, `
			INSERT INTO criteria (id, node_id, created_at, updated_at, description, is_completed)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, c.ID, c.NodeID, c.CreatedAT, c.UpdatedAt, c.Description, c.IsCompleted)
		if err != nil {
			return fmt.Errorf("save criterion: %w", err)
		}
	}

	return tx.Commit(ctx)
}
