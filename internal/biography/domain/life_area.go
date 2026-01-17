package domain

import (
	common "go-interview/internal/common/domain"
)

type LifeArea struct {
	common.UpdatableEntity
	Title common.Title
	Goal  common.Goal

	Criteria []*Criterion
}

func NewLifeArea(title common.Title, goal common.Goal) LifeArea {
	lifeArea := LifeArea{
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
		la.Criteria = append(la.Criteria, &criterion)
	}
}
