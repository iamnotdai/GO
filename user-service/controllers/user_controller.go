package controllers

import (
	"net/http"
	"user-service/dtos"
	"user-service/services"

	"github.com/gin-gonic/gin"
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
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := u.UserService.CreateUser(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// FetchUsers godoc
// @Summary Get list of users
// @Tags Users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (u *UserController) FetchUsers(c *gin.Context) {
	users, err := u.UserService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
	user, err := u.UserService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
