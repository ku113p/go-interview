package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	listlifeareas "go-interview/internal/biography/app/queries/list_life_areas"
)

// ListLifeAreasHandlerHTTP handles GET /life-areas requests.
type ListLifeAreasHandlerHTTP struct {
	useCase *listlifeareas.ListLifeAreaHandler
}

func NewListLifeAreasHandlerHTTP(uc *listlifeareas.ListLifeAreaHandler) *ListLifeAreasHandlerHTTP {
	return &ListLifeAreasHandlerHTTP{useCase: uc}
}

func (h *ListLifeAreasHandlerHTTP) Handle(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
		return
	}

	query := listlifeareas.ListLifeAreaQuery{UserID: userID}
	result, err := h.useCase.Handle(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, result)
}
