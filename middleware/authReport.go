package middleware

import (
	"net/http"

	"moodleinix/constant"

	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != constant.API_TOKEN {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}