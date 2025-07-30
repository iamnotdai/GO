package controllers

import (
	"net/http"
	"user-service/dtos"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService}
}

// CreateUser godoc
// @Summary Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param input body dtos.CreateUserInput true "User info"
// @Success 201 {object} gin.H{"message": string}
// @Failure 400 {object} gin.H{"error": string}
// @Router /users [post]
func (u *UserController) CreateUser(c *gin.Context) {
	var input dtos.CreateUserInput

	log.Info().
		Str("action", "CreateUser").
		Str("email", input.Email).
		Msg("Received create user request")

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Warn().
			Err(err).
			Str("action", "CreateUser").
			Msg("Invalid input format")

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := u.UserService.CreateUser(input); err != nil {
		log.Warn().
			Err(err).
			Str("action", "CreateUser").
			Str("email", input.Email).
			Msg("Failed to create user")

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info().
		Str("action", "CreateUser").
		Str("email", input.Email).
		Msg("User created successfully")

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// FetchUsers godoc
// @Summary Get list of users
// @Tags Users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (u *UserController) FetchUsers(c *gin.Context) {
	log.Info().
		Str("action", "FetchUsers").
		Msg("Fetching all users")

	users, err := u.UserService.GetAllUsers()
	if err != nil {
		log.Error().
			Err(err).
			Str("action", "FetchUsers").
			Msg("Failed to fetch users")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().
		Str("action", "FetchUsers").
		Int("count", len(users)).
		Msg("Fetched users successfully")

	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} gin.H{"error": string}
// @Router /users/{id} [get]
func (u *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	log.Info().
		Str("action", "GetUserByID").
		Str("user_id", id).
		Msg("Fetching user by ID")

	user, err := u.UserService.GetUserByID(id)
	if err != nil {
		log.Warn().
			Err(err).
			Str("action", "GetUserByID").
			Str("user_id", id).
			Msg("User not found")

		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Info().
		Str("action", "GetUserByID").
		Str("user_id", id).
		Msg("Fetched user successfully")

	c.JSON(http.StatusOK, user)
}
