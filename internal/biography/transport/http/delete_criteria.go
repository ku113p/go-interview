package http

import (
	"encoding/json"
	"errors"
	"net/http"

	usecase "go-interview/internal/biography/app/commands/delete_criteria"
	"go-interview/internal/biography/domain"
)

type DeleteCriteriaHandlerHTTP struct {
	useCase *usecase.DeleteCriteriaHandler
}

func NewDeleteCriteriaHandlerHTTP(uc *usecase.DeleteCriteriaHandler) *DeleteCriteriaHandlerHTTP {
	return &DeleteCriteriaHandlerHTTP{
		useCase: uc,
	}
}

func (h *DeleteCriteriaHandlerHTTP) Handle(w http.ResponseWriter, r *http.Request) {
	var cmd usecase.DeleteCriteriaCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	_, err := h.useCase.Handle(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			http.Error(w, "forbidden", http.StatusForbidden)
		} else if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "criteria not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
