package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"team-service/handlers"
	"team-service/middlewares"
	"team-service/repositories"
	"team-service/services"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	teamRepo := repositories.TeamRepository(db)
	teamService := services.TeamService(teamRepo)
	teamHandler := handlers.TeamHandler(teamService)

	// Public routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/")

	// Authenticated routes
	api.Use(middlewares.AuthMiddleware())

	// Only managers can create or modify teams
	api.Use(middlewares.RequireManagerRole())

	teamHandler.RegisterRoutes(api)

	return router
}
