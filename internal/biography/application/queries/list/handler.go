package list

import (
	"context"
	"go-interview/internal/biography/domain"

	"github.com/google/uuid"
)

type Handler struct {
	repo domain.LifeAreaRepository
}

func NewHandler(repo domain.LifeAreaRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, query Query) (Result, error) {
	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return Result{}, err
	}

	las, err := h.repo.ListLifeAreas(ctx, userID)
	if err != nil {
		return Result{}, err
	}
	result := make([]*AreaLife, 0, len(las))
	for _, la := range las {
		var parentID *string
		if la.ParentID != nil {
			idStr := la.ParentID.String()
			parentID = &idStr
		}
		result = append(result, &AreaLife{
			ID:        la.ID.String(),
			Title:     la.Title.String(),
			Goal:      la.Goal.String(),
			ParentID:  parentID,
			CreatedAt: la.CreatedAt.String(),
			UpdatedAt: la.UpdatedAt.String(),
		})
	}
	return Result{
		Items: result,
	}, nil
}
