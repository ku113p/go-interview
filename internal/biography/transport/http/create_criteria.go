package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	usecase "go-interview/internal/biography/app/commands/create_criteria"
	"go-interview/internal/biography/domain"
)

// CreateCriteriaHandlerHTTP handles POST /life-areas/:id/criteria requests.
type CreateCriteriaHandlerHTTP struct {
	useCase *usecase.CreateCriteriaHandler
}

func NewCreateCriteriaHandlerHTTP(uc *usecase.CreateCriteriaHandler) *CreateCriteriaHandlerHTTP {
	return &CreateCriteriaHandlerHTTP{useCase: uc}
}

func (h *CreateCriteriaHandlerHTTP) Handle(c *gin.Context) {
	lifeAreaID := c.Param("id")
	if lifeAreaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "life_area_id path parameter is required"})
		return
	}

	var cmd usecase.CreateCriteriaCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cmd.LifeAreaID = lifeAreaID

	result, err := h.useCase.Handle(c.Request.Context(), cmd)
	if err != nil {
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

	c.JSON(http.StatusCreated, result)
}
