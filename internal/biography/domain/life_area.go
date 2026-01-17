package domain

import (
	common "go-interview/internal/common/domain"

	"github.com/google/uuid"
)

type LifeArea struct {
	common.UpdatableEntity

	UserID   uuid.UUID
	ParentID *uuid.UUID

	Title common.Title
	Goal  common.Goal

	Criteria []*Criterion
}

func NewLifeArea(userID uuid.UUID, parentID *uuid.UUID, title common.Title, goal common.Goal) *LifeArea {
	lifeArea := &LifeArea{
		UserID:   userID,
		ParentID: parentID,
		Title:    title,
		Goal:     goal,
		Criteria: make([]*Criterion, 0),
	}
	common.InitUpdatableEntity(&lifeArea.UpdatableEntity)
	return lifeArea
}

func (la *LifeArea) AddCriteria(descs ...common.Description) {
	for _, desc := range descs {
		criterion := NewCriterion(desc)
		la.Criteria = append(la.Criteria, criterion)
	}
}
