package postgres

import "go-interview/internal/biography/domain"

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
