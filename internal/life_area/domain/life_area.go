package domain

import (
	"time"

	"github.com/google/uuid"
)

type LifeArea struct {
	ID uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time

	UserID   uuid.UUID
	ParentID *uuid.UUID

	Title Title
	Goal  Goal
}

func NewLifeArea(id uuid.UUID, userID uuid.UUID, parentID *uuid.UUID, title Title, goal Goal) *LifeArea {
	now := time.Now()

	return &LifeArea{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
		ParentID:  parentID,
		Title:     title,
		Goal:      goal,
	}
}

func (la *LifeArea) ChangeParentID(parentID *uuid.UUID) {
	la.ParentID = parentID
	la.afterUpdate()
}

func (la *LifeArea) afterUpdate() {
	la.UpdatedAt = time.Now()
}

func (la *LifeArea) ChangeGoal(goal Goal) {
	la.Goal = goal
	la.afterUpdate()
}
