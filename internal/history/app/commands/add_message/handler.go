package addmessage

import (
	"context"
	"go-interview/internal/history/domain"

	"github.com/google/uuid"
)

type Repository interface {
	domain.UserHistoryUpdater
}

type AddMessageHandler struct {
	repo Repository
}

func NewAddMessageHandler(repo Repository) *AddMessageHandler {
	return &AddMessageHandler{
		repo: repo,
	}
}

func (h *AddMessageHandler) Handle(ctx context.Context, cmd AddMessageCommand) error {
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return err
	}

	return h.repo.AddMessageToHistory(ctx, userID, cmd.Message, cmd.Limit)
}
