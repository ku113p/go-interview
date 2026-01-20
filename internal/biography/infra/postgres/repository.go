package postgres

import (
	"context"
	"errors"
	"fmt"
	"go-interview/internal/biography/domain"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ domain.LifeAreaCreator = (*AreaRepository)(nil)
var _ domain.LifeAreaGetter = (*AreaRepository)(nil)
var _ domain.LifeAreaLister = (*AreaRepository)(nil)
var _ domain.LifeAreaDeleter = (*AreaRepository)(nil)
var _ domain.LifeAreaParentChanger = (*AreaRepository)(nil)
var _ domain.LifeAreaGoalChanger = (*AreaRepository)(nil)
var _ domain.CriteriaDeleter = (*AreaRepository)(nil)
var _ domain.CriteriaCreator = (*AreaRepository)(nil)
var _ domain.CriteriaNodeGetter = (*AreaRepository)(nil)

type AreaRepository struct {
	db *pgxpool.Pool
}

func NewAreaRepository(db *pgxpool.Pool) *AreaRepository {
	return &AreaRepository{
		db: db,
	}
}

func (r *AreaRepository) CreateLifeArea(ctx context.Context, la *domain.LifeArea) error {
	return r.save(ctx, la)
}

func (r *AreaRepository) save(ctx context.Context, area *domain.LifeArea) error {
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

	_, err = tx.Exec(ctx, "DELETE FROM criteria WHERE node_id = $1", lifeAreaSQL.ID)
	if err != nil {
		return fmt.Errorf("delete criteria: %w", err)
	}

	if len(area.Criteria) > 0 {
		values := make([]string, 0, len(area.Criteria))
		args := make([]any, 0, len(area.Criteria)*6)

		argPos := 1
		for _, c := range area.Criteria {
			values = append(values,
				fmt.Sprintf(
					"($%d,$%d,$%d,$%d,$%d,$%d)",
					argPos, argPos+1, argPos+2,
					argPos+3, argPos+4, argPos+5,
				),
			)

			c.NodeID = area.ID

			args = append(args,
				c.ID,
				c.NodeID,
				c.CreatedAt,
				c.UpdatedAt,
				c.Description.String(),
				c.IsCompleted,
			)

			argPos += 6
		}
		joinedValues := strings.Join(values, ",\n")
		query := fmt.Sprintf(
			"INSERT INTO criteria (id, node_id, created_at, updated_at, description, is_completed) VALUES %s",
			joinedValues,
		)

		if _, err = tx.Exec(ctx, query, args...); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *AreaRepository) CreateCriteria(ctx context.Context, criteria ...*domain.Criterion) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if len(criteria) == 0 {
		return nil
	}

	values := make([]string, 0, len(criteria))
	args := make([]any, 0, len(criteria)*6)

	argPos := 1
	for _, c := range criteria {
		values = append(values,
			fmt.Sprintf(
				"($%d,$%d,$%d,$%d,$%d,$%d)",
				argPos, argPos+1, argPos+2,
				argPos+3, argPos+4, argPos+5,
			),
		)

		args = append(args,
			c.ID,
			c.NodeID,
			c.CreatedAt,
			c.UpdatedAt,
			c.Description.String(),
			c.IsCompleted,
		)

		argPos += 6
	}

	joinedValues := strings.Join(values, ",\n")
	query := fmt.Sprintf(
		"INSERT INTO criteria (id, node_id, created_at, updated_at, description, is_completed) VALUES %s",
		joinedValues,
	)

	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *AreaRepository) GetLifeArea(ctx context.Context, id uuid.UUID) (*domain.LifeArea, error) {
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
			return nil, domain.ErrNotFound
		}
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
			&criterionSQL.CreatedAt,
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

func (r *AreaRepository) ListLifeAreas(ctx context.Context, userID uuid.UUID) ([]*domain.LifeArea, error) {
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

func (r *AreaRepository) DeleteLifeArea(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM nodes WHERE id = $1", id)
	return err
}

func (r *AreaRepository) ChangeParentID(ctx context.Context, ID uuid.UUID, parentID *uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE nodes
		SET parent_id = $1, updated_at = NOW()
		WHERE id = $2
	`, parentID, ID)
	return err
}

func (r *AreaRepository) ChangeGoal(ctx context.Context, ID uuid.UUID, goal domain.Goal) error {
	_, err := r.db.Exec(ctx, `
		UPDATE nodes
		SET goal = $1, updated_at = NOW()
		WHERE id = $2
	`, goal.String(), ID)
	return err
}

func (r *AreaRepository) DeleteCriteria(ctx context.Context, IDs ...uuid.UUID) error {
	inParams := make([]string, 0, len(IDs))
	args := make([]any, 0, len(IDs))
	for ind, ID := range IDs {
		inParams = append(
			inParams,
			fmt.Sprintf("$%d", ind+1),
		)
		args = append(args, ID)
	}
	joinedInParams := strings.Join(inParams, ",")

	query := fmt.Sprintf(
		"DELETE FROM criteria WHERE id IN (%s)",
		joinedInParams,
	)
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete criteria: %w", err)
	}

	return err
}

func (r *AreaRepository) GetCriteriaNodeIDs(ctx context.Context, IDs ...uuid.UUID) (map[uuid.UUID]uuid.UUID, error) {
	result := make(map[uuid.UUID]uuid.UUID, len(IDs))
	if len(IDs) == 0 {
		return result, nil
	}

	inParams := make([]string, 0, len(IDs))
	args := make([]any, 0, len(IDs))
	for ind, ID := range IDs {
		inParams = append(inParams, fmt.Sprintf("$%d", ind+1))
		args = append(args, ID)
	}
	joinedInParams := strings.Join(inParams, ",")

	query := fmt.Sprintf(
		"SELECT id, node_id FROM criteria WHERE id IN (%s)",
		joinedInParams,
	)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get criteria nodes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var criterionID uuid.UUID
		var nodeID uuid.UUID
		if err := rows.Scan(&criterionID, &nodeID); err != nil {
			return nil, fmt.Errorf("scan criteria node: %w", err)
		}
		result[criterionID] = nodeID
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate criteria nodes: %w", err)
	}

	if len(result) != len(IDs) {
		return nil, domain.ErrNotFound
	}

	return result, nil
}
