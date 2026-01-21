package postgres

import "go-interview/internal/memory/domain"

func factToSQL(fact *domain.Fact) *FactSQL {
	return &FactSQL{
		ID:        fact.ID,
		CreatedAt: fact.CreatedAt,
		NodeID:    fact.NodeID,
		Info:      fact.Info.String(),
		DateTime:  fact.DateTime,
	}
}

func (dto *FactSQL) ToDomain() *domain.Fact {
	return &domain.Fact{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		NodeID:    dto.NodeID,
		Info:      domain.NewInfo(dto.Info),
		DateTime:  dto.DateTime,
	}
}

func embeddingToSQL(embedding *domain.Embedding) *EmbeddingSQL {
	return &EmbeddingSQL{
		ID:        embedding.ID,
		CreatedAt: embedding.CreatedAt,
		NodeID:    embedding.NodeID,
		Vector:    []float64(embedding.Vector),
		Content:   embedding.Content.String(),
	}
}

func (dto *EmbeddingSQL) ToDomain() *domain.Embedding {
	return &domain.Embedding{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		NodeID:    dto.NodeID,
		Vector:    domain.NewVector(dto.Vector),
		Content:   domain.NewContent(dto.Content),
	}
}
