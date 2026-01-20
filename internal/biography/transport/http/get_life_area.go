package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	usecase "go-interview/internal/biography/app/queries/get_life_area"
	"go-interview/internal/biography/domain"
)

// GetLifeAreaHandlerHTTP handles GET /life-areas/:id requests.
type GetLifeAreaHandlerHTTP struct {
	useCase *usecase.GetLifeAreaHandler
}

func NewGetLifeAreaHandlerHTTP(uc *usecase.GetLifeAreaHandler) *GetLifeAreaHandlerHTTP {
	return &GetLifeAreaHandlerHTTP{useCase: uc}
}

func (h *GetLifeAreaHandlerHTTP) Handle(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not found in path"})
		return
	}

	query := usecase.GetLifeAreaQuery{ID: id}
	result, err := h.useCase.Handle(c.Request.Context(), query)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "life area not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, result)
}
