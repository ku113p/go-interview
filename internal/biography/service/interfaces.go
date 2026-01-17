package service

import (
	"context"
	"go-interview/internal/biography/domain"
	common "go-interview/internal/common/domain"
)

type CriteriaGenerator interface {
	Generate(ctx context.Context, area *domain.LifeArea) ([]common.Description, error)
}
