package postgres

import (
	"go-interview/internal/criterion/domain"
)

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
