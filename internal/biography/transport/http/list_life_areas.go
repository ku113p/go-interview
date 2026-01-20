package http

import (
	"encoding/json"
	"net/http"

	usecase "go-interview/internal/biography/app/queries/list_life_areas"
)

type ListLifeAreasHandlerHTTP struct {
	useCase *usecase.ListLifeAreaHandler
}

func NewListLifeAreasHandlerHTTP(uc *usecase.ListLifeAreaHandler) *ListLifeAreasHandlerHTTP {
	return &ListLifeAreasHandlerHTTP{
		useCase: uc,
	}
}

func (h *ListLifeAreasHandlerHTTP) Handle(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	query := usecase.ListLifeAreaQuery{
		UserID: userID,
	}

	result, err := h.useCase.Handle(r.Context(), query)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
