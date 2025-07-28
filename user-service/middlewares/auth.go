package middlewares

import (
	"net/http"
	"strings"
	"user-service/services"
	"user-service/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Xác thực user từ DB
		user, err := userService.GetUserByID(claims.UserID)
		if err != nil || user.Role != claims.Role {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user or role mismatch"})
			c.Abort()
			return
		}

		// Set userID and role to context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
