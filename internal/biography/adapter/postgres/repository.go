package postgres

import (
	"context"
	"fmt"
	"go-interview/internal/biography/domain"

	"github.com/google/uuid"
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

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.LifeArea, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var lifeAreaSQL LifeAreaSQL
	if err := tx.QueryRow(ctx, `
		SELECT * FROM nodes WHERE id = $1
	`, id).Scan(
		&lifeAreaSQL.ID,
		&lifeAreaSQL.ParentID,
		&lifeAreaSQL.UserID,
		&lifeAreaSQL.CreatedAt,
		&lifeAreaSQL.UpdatedAt,
		&lifeAreaSQL.Title,
		&lifeAreaSQL.Goal,
	); err != nil {
		return nil, fmt.Errorf("get node: %w", err)
	}
	area := lifeAreaSQL.ToDomain()
	rows, err := tx.Query(ctx, `
		SELECT * FROM criteria WHERE node_id = $1
	`, id)
	if err != nil {
		return nil, fmt.Errorf("get criteria: %w", err)
	}
	for rows.Next() {
		var criterionSQL CriterionSQL
		if err := rows.Scan(
			&criterionSQL.ID,
			&criterionSQL.NodeID,
			&criterionSQL.CreatedAT,
			&criterionSQL.UpdatedAt,
			&criterionSQL.Description,
			&criterionSQL.IsCompleted,
		); err != nil {
			return nil, fmt.Errorf("scan criterion: %w", err)
		}
		area.Criteria = append(area.Criteria, criterionSQL.ToDomain())
	}

	return area, nil
}
