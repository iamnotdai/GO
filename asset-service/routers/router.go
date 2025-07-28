package routers

import (
	"asset-service/controllers"
	"asset-service/middlewares"
	"asset-service/repositories"
	"asset-service/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Repository
	folderRepo := repositories.NewFolderRepository(db)
	noteRepo := repositories.NewNoteRepository(db)

	// Service
	folderService := services.NewFolderService(folderRepo, noteRepo)
	noteService := services.NewNoteService(noteRepo)

	// Controller
	folderController := &controllers.FolderController{
		FolderService: folderService,
	}
	noteController := &controllers.NoteController{
		NoteService: noteService,
	}

	// Authenticated routes
	auth := router.Group("/")
	auth.Use(middlewares.AuthMiddleware())

	// Folder routes
	folder := router.Group("/folders")
	{
		folder.POST("", folderController.CreateFolder)                    // POST /api/folders
		folder.GET("/:id", folderController.GetFolder)                    // GET  /api/folders/:id
		folder.POST("/:id/share", folderController.ShareFolder)           // POST /api/folders/:id/share
		folder.DELETE("/:id/share/:userId", folderController.RevokeShare) // DELETE /api/folders/:id/share/:userId
	}

	// Note routes
	note := router.Group("/notes")
	{
		note.POST("", noteController.CreateNote)                      // POST /api/notes
		note.GET("/:id", noteController.GetNote)                      // GET  /api/notes/:id
		note.POST("/:id/share", noteController.ShareNote)             // POST /api/notes/:id/share
		note.DELETE("/:id/share/:userId", noteController.RevokeShare) // DELETE /api/notes/:id/share/:userId
	}

	return router
}
