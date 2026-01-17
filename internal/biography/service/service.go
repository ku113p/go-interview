package service

import (
	"context"
	"fmt"
	"go-interview/internal/biography/domain"
	common "go-interview/internal/common/domain"

	"github.com/google/uuid"
)

type BiographyService struct {
	repo domain.LifeAreaRepository
	ai   CriteriaGenerator
}

func NewBiographyService(r domain.LifeAreaRepository, g CriteriaGenerator) *BiographyService {
	return &BiographyService{
		repo: r,
		ai:   g,
	}
}

func (s *BiographyService) CreateLifeArea(ctx context.Context, userID uuid.UUID, parentID *uuid.UUID, title common.Title, goal common.Goal) error {
	area := domain.NewLifeArea(userID, parentID, title, goal)

	criteriaDescriptions, err := s.ai.Generate(ctx, area)
	if err != nil {
		return fmt.Errorf("failed to generate criteria: %w", err)
	}

	area.AddCriteria(criteriaDescriptions...)

	if err := s.repo.Save(ctx, area); err != nil {
		return fmt.Errorf("failed to save life area: %w", err)
	}

	return nil
}
