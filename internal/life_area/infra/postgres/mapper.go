package postgres

import (
	"go-interview/internal/life_area/domain"
)

func lifeAreaToSQL(area *domain.LifeArea) *LifeAreaSQL {
	return &LifeAreaSQL{
		ID:        area.ID,
		ParentID:  area.ParentID,
		UserID:    area.UserID,
		CreatedAt: area.CreatedAt,
		UpdatedAt: area.UpdatedAt,
		Title:     area.Title.String(),
		Goal:      area.Goal.String(),
	}
}

func (dto *LifeAreaSQL) ToDomain() *domain.LifeArea {
	return &domain.LifeArea{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		ParentID:  dto.ParentID,
		UserID:    dto.UserID,
		Title:     domain.NewTitle(dto.Title),
		Goal:      domain.NewGoal(dto.Goal),
	}
}
