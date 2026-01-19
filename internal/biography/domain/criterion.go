package domain

import (
	common "go-interview/internal/common/domain"

	"github.com/google/uuid"
)

type Criterion struct {
	common.UpdatableEntity

	NodeID      uuid.UUID
	Description Description
	IsCompleted bool
}

func NewCriterion(description Description) *Criterion {
	criterion := &Criterion{
		Description: description,
		IsCompleted: false,
	}
	common.InitUpdatableEntity(&criterion.UpdatableEntity)
	return criterion
}
