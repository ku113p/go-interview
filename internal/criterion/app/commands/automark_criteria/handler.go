package automarkcriteria

import (
	"context"
	"go-interview/internal/criterion/domain"

	"github.com/google/uuid"
)

type Repository interface {
	domain.CriteriaLister
}

type AiSerivce interface {
	domain.AICriteriaMeeter
}

type AutomarkCriteriaHandler struct {
	repo Repository
	ai   AiSerivce
}

func NewAutomarkCriteriaHandler(repo Repository, ai AiSerivce) *AutomarkCriteriaHandler {
	return &AutomarkCriteriaHandler{
		repo: repo,
		ai:   ai,
	}
}

func (h *AutomarkCriteriaHandler) Handle(ctx context.Context, cmd AutomarkCriteriaCommand) (*AutomarkCriteriaResult, error) {
	nodeID, err := uuid.Parse(cmd.NodeID)
	if err != nil {
		return nil, err
	}

	criteria, err := h.repo.GetCriteriaByNode(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	allMarked := true
	for _, c := range criteria {
		if !c.IsCompleted {
			allMarked = false
			break
		}
	}
	if allMarked {
		return &AutomarkCriteriaResult{
			AllMarked: true,
		}, nil
	}

	meets, err := h.ai.MeetCriteria(ctx, cmd.Text, criteria)
	if err != nil {
		return nil, err
	}
	allMarked = true
	for _, m := range meets {
		if !m {
			allMarked = false
			break
		}
	}

	return &AutomarkCriteriaResult{
		AllMarked: allMarked,
	}, nil
}
