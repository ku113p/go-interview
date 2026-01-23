package postgres

import (
	"context"
	"fmt"
	"go-interview/internal/criterion/domain"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ domain.CriteriaDeleter = (*PostgresRepository)(nil)
var _ domain.CriteriaCreator = (*PostgresRepository)(nil)
var _ domain.CriteriaLister = (*PostgresRepository)(nil)
var _ domain.CriteriaMarker = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) CreateCriteria(ctx context.Context, criteria ...*domain.Criterion) error {
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

func (r *PostgresRepository) DeleteCriteria(ctx context.Context, IDs ...uuid.UUID) error {
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

func (r *PostgresRepository) GetCriteriaByNode(ctx context.Context, nodeID uuid.UUID) ([]*domain.Criterion, error) {
	query := `
		SELECT * FROM criteria WHERE node_id = $1
	`

	rows, err := r.db.Query(ctx, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("get criteria: %w", err)
	}
	defer rows.Close()

	criteria := make([]*domain.Criterion, 0)
	for rows.Next() {
		var criterionSQL CriterionSQL
		if err := criterionSQL.Scan(rows); err != nil {
			return nil, fmt.Errorf("scan criterion: %w", err)
		}
		criteria = append(criteria, criterionSQL.ToDomain())
	}

	return criteria, nil
}

func (r *PostgresRepository) MarkCriteria(ctx context.Context, ids ...uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	query := `
		UPDATE criteria
		SET is_completed = true
		WHERE id = ANY($1)
	`

	_, err := r.db.Exec(ctx, query, ids)
	return err
}
