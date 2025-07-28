package controllers

import (
	"asset-service/models"
	"asset-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FolderController struct {
	FolderService *services.FolderService
}

// Tạo folder mới
func (fc *FolderController) CreateFolder(c *gin.Context) {
	var folder models.Folder
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Lấy owner_id từ JWT middleware (giả sử key là "user_id")
	userID, _ := c.Get("user_id")
	folder.ID = uuid.New()
	folder.OwnerID = userID.(uuid.UUID)

	if err := fc.FolderService.CreateFolder(&folder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create folder"})
		return
	}

	c.JSON(http.StatusCreated, folder)
}

// Lấy folder theo ID nếu có quyền
func (fc *FolderController) GetFolder(c *gin.Context) {
	idParam := c.Param("id")
	folderID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	folder, err := fc.FolderService.GetFolderByID(folderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}

	_, ok := fc.FolderService.HasAccess(folder, userID)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "no access to this folder"})
		return
	}

	c.JSON(http.StatusOK, folder)
}

// Chia sẻ folder với user khác
func (fc *FolderController) ShareFolder(c *gin.Context) {
	var input struct {
		UserID string `json:"user_id"`
		Access string `json:"access"` // "read" hoặc "write"
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id"})
		return
	}

	targetUserID, err := uuid.Parse(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := fc.FolderService.ShareFolderWithUser(folderID, targetUserID, input.Access); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to share folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "folder shared successfully"})
}

// Thu hồi chia sẻ
func (fc *FolderController) RevokeShare(c *gin.Context) {
	folderID, _ := uuid.Parse(c.Param("id"))
	userID, _ := uuid.Parse(c.Param("userId"))

	if err := fc.FolderService.RevokeFolderShare(folderID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke share"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "access revoked"})
}
