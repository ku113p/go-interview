package http

import "github.com/gin-gonic/gin"

func getPathParam(c *gin.Context, key string) (string, bool) {
	value := c.Param(key)
	if value == "" {
		return "", false
	}
	return value, true
}
