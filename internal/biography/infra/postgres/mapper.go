package postgres

import (
	"go-interview/internal/biography/domain"
)

func lifeAreaToSQL(area *domain.LifeArea) *LifeAreaSQL {
	criteria := make([]*CriterionSQL, 0, len(area.Criteria))
	for _, c := range area.Criteria {
		criteria = append(
			criteria,
			criterionToSQL(c),
		)
	}

	node := &LifeAreaSQL{
		ID:        area.ID,
		ParentID:  area.ParentID,
		UserID:    area.UserID,
		CreatedAt: area.CreatedAt,
		UpdatedAt: area.UpdatedAt,
		Title:     area.Title.String(),
		Goal:      area.Goal.String(),
	}

	return node
}

func criterionToSQL(criterion *domain.Criterion) *CriterionSQL {
	return &CriterionSQL{
		ID:          criterion.ID,
		NodeID:      criterion.NodeID,
		CreatedAt:   criterion.CreatedAt,
		UpdatedAt:   criterion.UpdatedAt,
		Description: criterion.Description.String(),
		IsCompleted: criterion.IsCompleted,
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
		Criteria:  make([]*domain.Criterion, 0),
	}
}

func (dto *CriterionSQL) ToDomain() *domain.Criterion {
	return &domain.Criterion{
		ID:          dto.ID,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		NodeID:      dto.NodeID,
		Description: domain.NewDescription(dto.Description),
		IsCompleted: dto.IsCompleted,
	}
}
