package domain

import (
	common "go-interview/internal/common/domain"
)

type Criterion struct {
	common.UpdatableEntity
	Description common.Description
	IsCompleted bool
}

func NewCriterion(description common.Description) *Criterion {
	criterion := &Criterion{
		Description: description,
		IsCompleted: false,
	}
	common.InitUpdatableEntity(&criterion.UpdatableEntity)
	return criterion
}
