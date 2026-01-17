package postgres

import (
	"go-interview/internal/biography/domain"
	common "go-interview/internal/common/domain"
)

func toSQL(area *domain.LifeArea) (*LifeAreaSQL, []*CriterionSQL) {
	node := &LifeAreaSQL{
		ID:        area.ID,
		ParentID:  *area.ParentID,
		UserID:    area.UserID,
		CreatedAt: area.CreatedAt,
		UpdatedAt: area.UpdatedAt,
		Title:     area.Title.String(),
		Goal:      area.Goal.String(),
	}

	criteria := make([]*CriterionSQL, 0, len(area.Criteria))
	for _, c := range area.Criteria {
		criteria = append(criteria, &CriterionSQL{
			ID:          c.ID,
			NodeID:      area.ID,
			CreatedAT:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
			Description: c.Description.String(),
			IsCompleted: c.IsCompleted,
		})
	}

	return node, criteria
}

func (dto *LifeAreaSQL) ToDomain() *domain.LifeArea {
	return &domain.LifeArea{
		UpdatableEntity: common.UpdatableEntity{
			Entity: common.Entity{
				ID:        dto.ID,
				CreatedAt: dto.CreatedAt,
			},
			UpdatedAt: dto.UpdatedAt,
		},
		ParentID: &dto.ParentID,
		UserID:   dto.UserID,
		Title:    common.NewTitle(dto.Title),
		Goal:     common.NewGoal(dto.Goal),
		Criteria: make([]*domain.Criterion, 0),
	}
}

func (dto *CriterionSQL) ToDomain() *domain.Criterion {
	return &domain.Criterion{
		UpdatableEntity: common.UpdatableEntity{
			Entity: common.Entity{
				ID:        dto.ID,
				CreatedAt: dto.CreatedAT,
			},
			UpdatedAt: dto.UpdatedAt,
		},
		Description: common.NewDescription(dto.Description),
		IsCompleted: dto.IsCompleted,
	}
}
