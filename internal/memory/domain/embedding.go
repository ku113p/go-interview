package domain

import "go-interview/internal/common/domain"

type Embedding struct {
	domain.Entity
	Vector  Vector
	Content domain.Content
}

func NewEmbedding(vector Vector, content domain.Content) *Embedding {
	embedding := &Embedding{
		Vector:  vector,
		Content: content,
	}
	domain.InitEntity(&embedding.Entity)
	return embedding
}
