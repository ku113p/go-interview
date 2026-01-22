package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"go-interview/internal/user/domain"
)

var _ domain.UserCreator = (*Repository)(nil)
var _ domain.UserByExternalIDGetter = (*Repository)(nil)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	userSQL := userToSQL(user)

	_, err := r.db.Exec(ctx, `
		INSERT INTO users (id, created_at, external_id)
		VALUES ($1, $2, $3)
	`, userSQL.ID, userSQL.CreatedAt, userSQL.ExternalID)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *Repository) GetUserByExternalID(ctx context.Context, externalID string) (*domain.User, error) {
	var userSQL UserSQL
	row := r.db.QueryRow(ctx, `
		SELECT id, created_at, external_id
		FROM users
		WHERE external_id = $1
	`, externalID)
	if err := userSQL.Scan(row); err != nil {
		return nil, fmt.Errorf("get user by external id: %w", err)
	}

	return userSQL.ToDomain(), nil
}
