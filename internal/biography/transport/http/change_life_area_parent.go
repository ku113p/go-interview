package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	usecase "go-interview/internal/biography/app/commands/change_life_area_parent"
	"go-interview/internal/biography/domain"
)

// ChangeLifeAreaParentHandlerHTTP handles PATCH /life-areas/:id/parent requests.
type ChangeLifeAreaParentHandlerHTTP struct {
	useCase *usecase.ChangeLifeAreaParentHandler
}

func NewChangeLifeAreaParentHandlerHTTP(uc *usecase.ChangeLifeAreaParentHandler) *ChangeLifeAreaParentHandlerHTTP {
	return &ChangeLifeAreaParentHandlerHTTP{useCase: uc}
}

func (h *ChangeLifeAreaParentHandlerHTTP) Handle(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID path parameter is required"})
		return
	}

	var cmd usecase.ChangeLifeAreaParentCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cmd.ID = id

	if _, err := h.useCase.Handle(c.Request.Context(), cmd); err != nil {
		switch {
		case errors.Is(err, domain.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "life area not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
