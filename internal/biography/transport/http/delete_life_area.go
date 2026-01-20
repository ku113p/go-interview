package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	usecase "go-interview/internal/biography/app/commands/delete_life_area"
	"go-interview/internal/biography/domain"
)

// DeleteLifeAreaHandlerHTTP handles DELETE /life-areas/:id requests.
type DeleteLifeAreaHandlerHTTP struct {
	useCase *usecase.DeleteLifeAreaHandler
}

func NewDeleteLifeAreaHandlerHTTP(uc *usecase.DeleteLifeAreaHandler) *DeleteLifeAreaHandlerHTTP {
	return &DeleteLifeAreaHandlerHTTP{useCase: uc}
}

func (h *DeleteLifeAreaHandlerHTTP) Handle(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID path parameter is required"})
		return
	}

	var cmd usecase.DeleteLifeAreaCommand
	_ = json.NewDecoder(c.Request.Body).Decode(&cmd) // optional user_id for future use
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
