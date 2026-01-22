package createuser

import (
	"context"
	"go-interview/internal/user/domain"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, user *domain.User) error
}

type IDGenerator interface {
	Generate() (uuid.UUID, error)
}

type Handler struct {
	repo  Repository
	genID IDGenerator
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd CreateUserCommand) (*CreateUserResult, error) {
	id, err := h.genID.Generate()
	if err != nil {
		return nil, err
	}

	user := domain.NewUser(id, cmd.ExternalID)

	err = h.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &CreateUserResult{
		ID: user.ID.String(),
	}, nil
}
