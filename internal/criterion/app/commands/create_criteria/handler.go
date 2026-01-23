package create_criteria

import (
	"context"
	"go-interview/internal/criterion/domain"

	"github.com/google/uuid"
)

type CreateCriteriaRepository interface {
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
	nodeID, err := uuid.Parse(cmd.NodeID)
	if err != nil {
		return nil, err
	}

	criteriaToCreate := make([]*domain.Criterion, 0, len(cmd.Descriptions))
	criteriaIDs := make([]string, 0, len(cmd.Descriptions))
	for _, d := range cmd.Descriptions {
		id, err := h.genID.Generate()
		if err != nil {
			return nil, err
		}
		criteriaIDs = append(criteriaIDs, id.String())

		criteriaToCreate = append(criteriaToCreate, domain.NewCriterion(
			id,
			nodeID,
			domain.NewDescription(d),
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
