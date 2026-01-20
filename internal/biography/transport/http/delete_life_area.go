package http

import (
	"encoding/json"
	"errors"
	"net/http"

	usecase "go-interview/internal/biography/app/commands/delete_life_area"
	"go-interview/internal/biography/domain"
)

type DeleteLifeAreaHandlerHTTP struct {
	useCase *usecase.DeleteLifeAreaHandler
}

func NewDeleteLifeAreaHandlerHTTP(uc *usecase.DeleteLifeAreaHandler) *DeleteLifeAreaHandlerHTTP {
	return &DeleteLifeAreaHandlerHTTP{
		useCase: uc,
	}
}

func (h *DeleteLifeAreaHandlerHTTP) Handle(w http.ResponseWriter, r *http.Request) {
	id, ok := getPathParam(r, "id")
	if !ok {
		http.Error(w, "ID path parameter is required", http.StatusBadRequest)
		return
	}

	var cmd usecase.DeleteLifeAreaCommand
	// DELETE request may not have a body. In a real-world app, user_id would come from auth.
	// For this exercise, we'll try to read it from the body.
	_ = json.NewDecoder(r.Body).Decode(&cmd)

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
