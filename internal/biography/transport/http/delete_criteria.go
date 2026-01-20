package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	usecase "go-interview/internal/biography/app/commands/delete_criteria"
	"go-interview/internal/biography/domain"
)

// DeleteCriteriaHandlerHTTP handles DELETE /criteria requests.
type DeleteCriteriaHandlerHTTP struct {
	useCase *usecase.DeleteCriteriaHandler
}

func NewDeleteCriteriaHandlerHTTP(uc *usecase.DeleteCriteriaHandler) *DeleteCriteriaHandlerHTTP {
	return &DeleteCriteriaHandlerHTTP{useCase: uc}
}

func (h *DeleteCriteriaHandlerHTTP) Handle(c *gin.Context) {
	var cmd usecase.DeleteCriteriaCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if _, err := h.useCase.Handle(c.Request.Context(), cmd); err != nil {
		switch {
		case errors.Is(err, domain.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "criteria not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
