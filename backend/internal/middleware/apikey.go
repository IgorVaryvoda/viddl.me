package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if apiKey == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "API key not configured"})
			c.Abort()
			return
		}

		providedKey := c.GetHeader("X-API-Key")
		if providedKey == "" {
			providedKey = c.Query("api_key")
		}

		if subtle.ConstantTimeCompare([]byte(providedKey), []byte(apiKey)) != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
