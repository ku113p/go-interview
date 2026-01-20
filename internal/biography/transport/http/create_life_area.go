package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	usecase "go-interview/internal/biography/app/commands/create_life_area"
)

// CreateLifeAreaHandlerHTTP handles POST /life-areas requests.
type CreateLifeAreaHandlerHTTP struct {
	useCase *usecase.CreateLifeAreaHandler
}

func NewCreateLifeAreaHandlerHTTP(uc *usecase.CreateLifeAreaHandler) *CreateLifeAreaHandlerHTTP {
	return &CreateLifeAreaHandlerHTTP{useCase: uc}
}

func (h *CreateLifeAreaHandlerHTTP) Handle(c *gin.Context) {
	var cmd usecase.CreateLifeAreaCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.useCase.Handle(c.Request.Context(), cmd)
	if err != nil {
		if strings.Contains(err.Error(), "invalid UUID") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, result)
}
