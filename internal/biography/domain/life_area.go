package domain

import (
	common "go-interview/internal/common/domain"

	"github.com/google/uuid"
)

type LifeArea struct {
	common.UpdatableEntity

	UserID   uuid.UUID
	ParentID *uuid.UUID

	Title Title
	Goal  Goal

	Criteria []*Criterion
}

func NewLifeArea(userID uuid.UUID, parentID *uuid.UUID, title Title, goal Goal) *LifeArea {
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

func (la *LifeArea) ChangeParentID(parentID *uuid.UUID) {
	la.ParentID = parentID
	la.Update()
}

func (la *LifeArea) ChangeGoal(goal Goal) {
	la.Goal = goal
	la.Update()
}

func (la *LifeArea) SetCriteria(descs ...Description) {
	la.Criteria = make([]*Criterion, 0, len(descs))
	for _, desc := range descs {
		criterion := NewCriterion(desc)
		la.Criteria = append(la.Criteria, criterion)
	}
	la.Update()
}
