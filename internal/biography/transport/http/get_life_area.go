package http

import (
	"encoding/json"
	"errors"
	"net/http"

	usecase "go-interview/internal/biography/app/queries/get_life_area"
	"go-interview/internal/biography/domain"
)

type GetLifeAreaHandlerHTTP struct {
	useCase *usecase.GetLifeAreaHandler
}

func NewGetLifeAreaHandlerHTTP(uc *usecase.GetLifeAreaHandler) *GetLifeAreaHandlerHTTP {
	return &GetLifeAreaHandlerHTTP{
		useCase: uc,
	}
}

func (h *GetLifeAreaHandlerHTTP) Handle(w http.ResponseWriter, r *http.Request) {
	params, ok := r.Context().Value(paramsContextKey).(map[string]string)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	id, ok := params["id"]
	if !ok {
		http.Error(w, "id not found in path", http.StatusBadRequest)
		return
	}

	query := usecase.GetLifeAreaQuery{
		ID: id,
	}

	result, err := h.useCase.Handle(r.Context(), query)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "life area not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}