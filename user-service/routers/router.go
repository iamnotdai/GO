package routers

import (
	"user-service/controllers"
	"user-service/middlewares"
	"user-service/repositories"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Repository
	userRepository := repositories.NewUserRepository(db)

	// Service
	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService(userRepository)

	// Controller
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)

	// user routes
	user := router.Group("/users")
	{
		user.POST("", userController.CreateUser)
		user.GET("", middlewares.AuthMiddleware(userService), userController.FetchUsers)
		user.GET("/:id", middlewares.AuthMiddleware(userService), userController.GetUserByID)
	}

	// auth routes
	authen := router.Group("/auth")
	{
		authen.POST("/login", authController.Login)
	}
	return router
}
