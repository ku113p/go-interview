package createfact

import (
	"context"
	"fmt"
	"go-interview/internal/memory/domain"
)

type Repository interface {
	domain.FactCreator
}

type Handler struct {
	repo  Repository
	genID domain.IDGenerator
}

func NewHandler(repo Repository, genID domain.IDGenerator) *Handler {
	return &Handler{
		repo:  repo,
		genID: genID,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	if cmd.Info == "" {
		return fmt.Errorf("info is required")
	}

	factID, err := h.genID.Generate()
	if err != nil {
		return fmt.Errorf("generate id: %w", err)
	}

	fact := domain.NewFact(
		factID,
		cmd.NodeID,
		domain.NewInfo(cmd.Info),
		cmd.DateTime,
	)

	if err := h.repo.CreateFact(ctx, fact); err != nil {
		return fmt.Errorf("create fact: %w", err)
	}

	return nil
}
