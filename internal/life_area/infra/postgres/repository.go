package postgres

import (
	"context"
	"errors"
	"fmt"
	"go-interview/internal/life_area/domain"

	common "go-interview/internal/common/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ domain.LifeAreaCreator = (*PostgresRepository)(nil)
var _ domain.LifeAreaGetter = (*PostgresRepository)(nil)
var _ domain.LifeAreaLister = (*PostgresRepository)(nil)
var _ domain.LifeAreaDeleter = (*PostgresRepository)(nil)
var _ domain.LifeAreaParentChanger = (*PostgresRepository)(nil)
var _ domain.LifeAreaGoalChanger = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) CreateLifeArea(ctx context.Context, la *domain.LifeArea) error {
	return r.save(ctx, la)
}

func (r *PostgresRepository) save(ctx context.Context, area *domain.LifeArea) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	lifeAreaSQL := lifeAreaToSQL(area)

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

	return tx.Commit(ctx)
}

func (r *PostgresRepository) GetLifeArea(ctx context.Context, id uuid.UUID) (*domain.LifeArea, error) {
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common.ErrNotFound
		}
		return nil, fmt.Errorf("get node: %w", err)
	}
	area := lifeAreaSQL.ToDomain()

	return area, nil
}

func (r *PostgresRepository) ListLifeAreas(ctx context.Context, userID uuid.UUID) ([]*domain.LifeArea, error) {
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

func (r *PostgresRepository) DeleteLifeArea(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM nodes WHERE id = $1", id)
	return err
}

func (r *PostgresRepository) ChangeParentID(ctx context.Context, ID uuid.UUID, parentID *uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE nodes
		SET parent_id = $1, updated_at = NOW()
		WHERE id = $2
	`, parentID, ID)
	return err
}

func (r *PostgresRepository) ChangeGoal(ctx context.Context, ID uuid.UUID, goal domain.Goal) error {
	_, err := r.db.Exec(ctx, `
		UPDATE nodes
		SET goal = $1, updated_at = NOW()
		WHERE id = $2
	`, goal.String(), ID)
	return err
}
