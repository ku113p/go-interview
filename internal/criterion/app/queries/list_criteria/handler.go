package listcriteria

import (
	"context"
	"go-interview/internal/criterion/domain"

	"github.com/google/uuid"
)

type Repository interface {
	domain.CriteriaLister
}

type ListCriteriaHandler struct {
	repo Repository
}

func NewListCriteriaHandler(repo Repository) *ListCriteriaHandler {
	return &ListCriteriaHandler{
		repo: repo,
	}
}

func (h *ListCriteriaHandler) Handle(ctx context.Context, query ListCriteriaQuery) (*ListCriteriaResult, error) {
	nodeID, err := uuid.Parse(query.NodeID)
	if err != nil {
		return nil, err
	}

	criteria, err := h.repo.GetCriteriaByNode(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	items := make([]Criterion, 0, len(criteria))
	for _, c := range criteria {
		items = append(items, Criterion{
			ID:          c.ID.String(),
			NodeID:      c.NodeID.String(),
			Description: c.Description.String(),
			IsCompleted: c.IsCompleted,
		})
	}

	return &ListCriteriaResult{
		Items: items,
	}, nil
}
