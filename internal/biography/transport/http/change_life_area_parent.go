package http

import (
	"encoding/json"
	"errors"
	"net/http"

	usecase "go-interview/internal/biography/app/commands/change_life_area_parent"
	"go-interview/internal/biography/domain"
)

type ChangeLifeAreaParentHandlerHTTP struct {
	useCase *usecase.ChangeLifeAreaParentHandler
}

func NewChangeLifeAreaParentHandlerHTTP(uc *usecase.ChangeLifeAreaParentHandler) *ChangeLifeAreaParentHandlerHTTP {
	return &ChangeLifeAreaParentHandlerHTTP{
		useCase: uc,
	}
}

func (h *ChangeLifeAreaParentHandlerHTTP) Handle(w http.ResponseWriter, r *http.Request) {
	id, ok := getPathParam(r, "id")
	if !ok {
		http.Error(w, "ID path parameter is required", http.StatusBadRequest)
		return
	}

	var cmd usecase.ChangeLifeAreaParentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cmd.ID = id

	_, err := h.useCase.Handle(r.Context(), cmd)
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

	w.WriteHeader(http.StatusNoContent)
}
