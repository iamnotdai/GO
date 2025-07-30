package middlewares

import (
	"net/http"
	"strings"
	"user-service/services"
	"user-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func AuthMiddleware(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Warn().
				Str("action", "AuthMiddleware").
				Str("ip", c.ClientIP()).
				Msg("Missing or invalid Authorization header")

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			log.Warn().
				Err(err).
				Str("action", "AuthMiddleware").
				Str("ip", c.ClientIP()).
				Msg("Failed to parse token")

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Xác thực user từ DB
		user, err := userService.GetUserByID(claims.UserID)
		if err != nil {
			log.Warn().
				Err(err).
				Str("action", "AuthMiddleware").
				Str("user_id", claims.UserID).
				Str("ip", c.ClientIP()).
				Msg("User not found during auth")

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user or role mismatch"})
			c.Abort()
			return
		}

		// Kiểm tra role
		if user.Role != claims.Role {
			log.Warn().
				Str("action", "AuthMiddleware").
				Str("user_id", user.ID.String()).
				Str("expected_role", claims.Role).
				Str("actual_role", user.Role).
				Str("ip", c.ClientIP()).
				Msg("Role mismatch during auth")

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user or role mismatch"})
			c.Abort()
			return
		}

		log.Info().
			Str("action", "AuthMiddleware").
			Str("user_id", user.ID.String()).
			Str("role", user.Role).
			Str("ip", c.ClientIP()).
			Msg("Authenticated successfully")

		// Set userID and role to context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
