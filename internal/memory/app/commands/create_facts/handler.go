package createfacts

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go-interview/internal/memory/domain"
)

type Repository interface {
	domain.FactCreator
}

type Handler struct {
	repo      Repository
	genID     domain.IDGenerator
	extractor domain.FactExtractor
}

func NewHandler(repo Repository, genID domain.IDGenerator, extractor domain.FactExtractor) *Handler {
	return &Handler{
		repo:      repo,
		genID:     genID,
		extractor: extractor,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	text := strings.TrimSpace(cmd.Text)
	goal := strings.TrimSpace(cmd.Goal)
	if text == "" {
		return errors.New("text is required")
	}
	if goal == "" {
		return errors.New("goal is required")
	}

	factInfos, err := h.extractor.ExtractFacts(ctx, text, goal)
	if err != nil {
		return fmt.Errorf("extract facts: %w", err)
	}
	if len(factInfos) == 0 {
		return errors.New("no facts extracted")
	}

	facts := make([]*domain.Fact, 0, len(factInfos))
	for _, info := range factInfos {
		if info == "" {
			continue
		}

		factID, err := h.genID.Generate()
		if err != nil {
			return fmt.Errorf("generate id: %w", err)
		}

		facts = append(facts, domain.NewFact(
			factID,
			cmd.NodeID,
			domain.NewInfo(info),
			cmd.DateTime,
		))
	}

	if len(facts) == 0 {
		return errors.New("no valid facts to create")
	}

	if err := h.repo.CreateFacts(ctx, facts); err != nil {
		return fmt.Errorf("create facts: %w", err)
	}

	return nil
}
