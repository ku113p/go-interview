package http

import (
	"encoding/json"
	"errors"
	"net/http"

	usecase "go-interview/internal/biography/app/commands/create_criteria"
	"go-interview/internal/biography/domain"
)

type CreateCriteriaHandlerHTTP struct {
	useCase *usecase.CreateCriteriaHandler
}

func NewCreateCriteriaHandlerHTTP(uc *usecase.CreateCriteriaHandler) *CreateCriteriaHandlerHTTP {
	return &CreateCriteriaHandlerHTTP{
		useCase: uc,
	}
}

func (h *CreateCriteriaHandlerHTTP) Handle(w http.ResponseWriter, r *http.Request) {
	lifeAreaID, ok := getPathParam(r, "id")
	if !ok {
		http.Error(w, "life_area_id path parameter is required", http.StatusBadRequest)
		return
	}

	var cmd usecase.CreateCriteriaCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cmd.LifeAreaID = lifeAreaID

	result, err := h.useCase.Handle(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			http.Error(w, "forbidden", http.StatusForbidden)
		} else if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "life area not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
