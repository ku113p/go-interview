package http

import (
	"encoding/json"
	"net/http"
	"strings"

	usecase "go-interview/internal/biography/app/commands/create_life_area"
)

type CreateLifeAreaHandlerHTTP struct {
	useCase *usecase.CreateLifeAreaHandler
}

func NewCreateLifeAreaHandlerHTTP(uc *usecase.CreateLifeAreaHandler) *CreateLifeAreaHandlerHTTP {
	return &CreateLifeAreaHandlerHTTP{
		useCase: uc,
	}
}

func (h *CreateLifeAreaHandlerHTTP) Handle(w http.ResponseWriter, r *http.Request) {
	var cmd usecase.CreateLifeAreaCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.useCase.Handle(r.Context(), cmd)
	if err != nil {
		if strings.Contains(err.Error(), "invalid UUID") {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
