package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-interview/internal/memory/domain"
)

var _ domain.FactCreator = (*Repository)(nil)
var _ domain.EmbeddingCreator = (*Repository)(nil)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateFact(ctx context.Context, fact *domain.Fact) error {
	factSQL := factToSQL(fact)

	_, err := r.db.Exec(ctx, `
		INSERT INTO facts (id, created_at, node_id, info, date_time)
		VALUES ($1, $2, $3, $4, $5)
	`, factSQL.ID, factSQL.CreatedAt, factSQL.NodeID, factSQL.Info, factSQL.DateTime)
	if err != nil {
		return fmt.Errorf("create fact: %w", err)
	}

	return nil
}

func (r *Repository) CreateEmbeddings(ctx context.Context, embeddings []*domain.Embedding) error {
	if len(embeddings) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, emb := range embeddings {
		embSQL := embeddingToSQL(emb)
		batch.Queue(`
			INSERT INTO embeddings (id, created_at, node_id, vector, content)
			VALUES ($1, $2, $3, $4, $5)
		`, embSQL.ID, embSQL.CreatedAt, embSQL.NodeID, embSQL.Vector, embSQL.Content)
	}

	br := r.db.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return fmt.Errorf("execute batch: %w", err)
	}

	return nil
}
