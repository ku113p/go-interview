package postgres

import (
	"context"
	"fmt"
	"go-interview/internal/biography/domain"
	"go-interview/internal/biography/port"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ port.AreaRepository = (*AreaRepository)(nil)

type AreaRepository struct {
	db *pgxpool.Pool
}

func NewAreaRepository(db *pgxpool.Pool) *AreaRepository {
	return &AreaRepository{
		db: db,
	}
}

func (r *AreaRepository) Save(ctx context.Context, area *domain.LifeArea) error {
	lifeAreaSQL, criterionSQL := toSQL(area)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO nodes (id, parent_id, user_id, created_at, updated_at, title, goal)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			parent_id = $2,
			user_id = $3,
			updated_at = $5,
			title = $6,
			goal = $7
	`, lifeAreaSQL.ID, lifeAreaSQL.ParentID, lifeAreaSQL.UserID, lifeAreaSQL.CreatedAt, lifeAreaSQL.UpdatedAt, lifeAreaSQL.Title, lifeAreaSQL.Goal)
	if err != nil {
		return fmt.Errorf("upsert node: %w", err)
	}

	_, err = tx.Exec(ctx, "DELETE FROM criteria WHERE node_id = $1", lifeAreaSQL.ID)
	if err != nil {
		return fmt.Errorf("delete criteria: %w", err)
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

func (r *AreaRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.LifeArea, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var lifeAreaSQL LifeAreaSQL
	row := tx.QueryRow(ctx, `
		SELECT * FROM nodes WHERE id = $1
	`, id)
	if err := lifeAreaSQL.Scan(row); err != nil {
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
		if err := criterionSQL.Scan(rows); err != nil {
			return nil, fmt.Errorf("scan criterion: %w", err)
		}
		area.Criteria = append(area.Criteria, criterionSQL.ToDomain())
	}

	return area, nil
}

func (r *AreaRepository) List(ctx context.Context, userID uuid.UUID) ([]*domain.LifeArea, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		SELECT * FROM nodes WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("get nodes: %w", err)
	}

	areas := make([]*domain.LifeArea, 0)
	for rows.Next() {
		var lifeAreaSQL LifeAreaSQL
		if err := lifeAreaSQL.Scan(rows); err != nil {
			return nil, fmt.Errorf("scan node: %w", err)
		}
		area := lifeAreaSQL.ToDomain()
		areas = append(areas, area)
	}

	return areas, nil
}

func (r *AreaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM nodes WHERE id = $1", id)
	return err
}
