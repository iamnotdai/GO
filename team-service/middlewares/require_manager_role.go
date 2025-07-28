package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireManagerRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing role"})
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok || role != "manager" {
			c.JSON(http.StatusForbidden, gin.H{"error": "manager role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
