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

func (la *LifeArea) ChangeParentID(parentID *uuid.UUID) {
	la.ParentID = parentID
	la.Update()
}

func (la *LifeArea) ChangeGoal(goal common.Goal) {
	la.Goal = goal
	la.Update()
}

func (la *LifeArea) SetCriteria(descs ...common.Description) {
	la.Criteria = make([]*Criterion, 0, len(descs))
	for _, desc := range descs {
		criterion := NewCriterion(desc)
		la.Criteria = append(la.Criteria, criterion)
	}
	la.Update()
}

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
