package controllers

import (
	"net/http"
	"user-service/dtos"
	"user-service/services"

	"github.com/gin-gonic/gin"
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

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := a.AuthService.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
