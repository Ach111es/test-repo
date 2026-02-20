package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	APIKeyHeader = "X-API-KEY"
)

// APIKeyMiddleware validates the API key from request header
func APIKeyMiddleware(validAPIKeys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(APIKeyHeader)

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "API key is required",
			})
			c.Abort()
			return
		}

		// Check if the provided API key is valid
		isValid := false
		for _, key := range validAPIKeys {
			if apiKey == key {
				isValid = true
				break
			}
		}

		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "Invalid API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
