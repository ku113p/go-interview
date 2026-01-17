package domain

import "go-interview/internal/common/domain"

type Fact struct {
	domain.UpdatableEntity
	Content domain.Content
}

func NewFact(content domain.Content) *Fact {
	fact := &Fact{
		Content: content,
	}
	domain.InitUpdatableEntity(&fact.UpdatableEntity)
	return fact
}
