package create_criteria

import (
	"context"
	"fmt"
	"go-interview/internal/biography/domain"

	"github.com/google/uuid"
)

type CreateCriteriaRepository interface {
	domain.LifeAreaGetter
	domain.CriteriaCreator
}

type CreateCriteriaHandler struct {
	repo  CreateCriteriaRepository
	genID domain.IDGenerator
}

func NewCreateCriteriaHandler(repo CreateCriteriaRepository, genID domain.IDGenerator) *CreateCriteriaHandler {
	return &CreateCriteriaHandler{
		repo:  repo,
		genID: genID,
	}
}

func (h *CreateCriteriaHandler) Handle(ctx context.Context, cmd CreateCriteriaCommand) (*CreateCriteriaResult, error) {
	lifeAreaID, err := uuid.Parse(cmd.LifeAreaID)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, err
	}

	lifeArea, err := h.repo.GetLifeArea(ctx, lifeAreaID)
	if err != nil {
		return nil, err
	}

	if lifeArea.UserID != userID {
		return nil, fmt.Errorf("user %s is not owner of life area %s", userID, lifeAreaID)
	}

	criteriaToCreate := make([]*domain.Criterion, 0, len(cmd.Criteria))
	criteriaIDs := make([]string, 0, len(cmd.Criteria))
	for _, c := range cmd.Criteria {
		id, err := h.genID.GenID()
		if err != nil {
			return nil, err
		}
		criteriaIDs = append(criteriaIDs, id.String())

		criteriaToCreate = append(criteriaToCreate, domain.NewCriterion(
			id,
			lifeAreaID,
			domain.NewDescription(c.Description),
		))
	}

	err = h.repo.CreateCriteria(ctx, criteriaToCreate...)
	if err != nil {
		return nil, err
	}

	return &CreateCriteriaResult{
		IDs: criteriaIDs,
	}, nil
}
