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

func (r *Repository) CreateFacts(ctx context.Context, facts []*domain.Fact) error {
	if len(facts) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, fact := range facts {
		fSQL := factToSQL(fact)
		batch.Queue(`
			INSERT INTO facts (id, created_at, node_id, info, date_time)
			VALUES ($1, $2, $3, $4, $5)
		`, fSQL.ID, fSQL.CreatedAt, fSQL.NodeID, fSQL.Info, fSQL.DateTime)
	}

	br := r.db.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return fmt.Errorf("execute batch: %w", err)
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
