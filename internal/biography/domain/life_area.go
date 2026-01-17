package domain

import (
	common "go-interview/internal/common/domain"
)

type LifeArea struct {
	common.UpdatableEntity
	Title common.Title
	Goal  common.Goal
}

func NewLifeArea(title common.Title, goal common.Goal) LifeArea {
	lifeArea := LifeArea{
		Title: title,
		Goal:  goal,
	}
	common.InitUpdatableEntity(&lifeArea.UpdatableEntity)
	return lifeArea
}
