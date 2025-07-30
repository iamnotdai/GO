package controllers

import (
	"net/http"
	"user-service/dtos"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService}
}

// Login godoc
// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body services.LoginInput true "Login info"
// @Success 200 {object} gin.H{"token": string}
// @Failure 400,401 {object} gin.H{"error": string}
// @Router /auth/login [post]
func (a *AuthController) Login(c *gin.Context) {
	var input dtos.LoginInput

	// Ghi log khi request đến
	log.Info().
		Str("action", "Login").
		Str("email", input.Email).
		Msg("Received login request")

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Warn().
			Err(err).
			Str("action", "Login").
			Msg("Invalid input format")

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := a.AuthService.Login(input.Email, input.Password)
	if err != nil {
		log.Warn().
			Err(err).
			Str("action", "Login").
			Str("email", input.Email).
			Msg("Login failed")

		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Info().
		Str("action", "Login").
		Str("email", input.Email).
		Msg("Login successful")

	c.JSON(http.StatusOK, gin.H{"token": token})
}
