package controllers

import (
	"asset-service/models"
	"asset-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NoteController struct {
	NoteService *services.NoteService
}

// Tạo ghi chú mới
func (nc *NoteController) CreateNote(c *gin.Context) {
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	note.ID = uuid.New()
	note.OwnerID = userID

	if err := nc.NoteService.CreateNote(&note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create note"})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// Lấy ghi chú theo ID
func (nc *NoteController) GetNote(c *gin.Context) {
	noteID, _ := uuid.Parse(c.Param("id"))
	userID := c.MustGet("user_id").(uuid.UUID)

	note, err := nc.NoteService.GetNoteByID(noteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	_, ok := nc.NoteService.HasAccess(note, userID)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "no access to this note"})
		return
	}

	c.JSON(http.StatusOK, note)
}

// Chia sẻ ghi chú
func (nc *NoteController) ShareNote(c *gin.Context) {
	var input struct {
		UserID string `json:"user_id"`
		Access string `json:"access"` // "read" hoặc "write"
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	noteID, _ := uuid.Parse(c.Param("id"))
	targetUserID, _ := uuid.Parse(input.UserID)

	if err := nc.NoteService.ShareNoteWithUser(noteID, targetUserID, input.Access); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to share note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "note shared successfully"})
}

// Thu hồi chia sẻ
func (nc *NoteController) RevokeShare(c *gin.Context) {
	noteID, _ := uuid.Parse(c.Param("id"))
	userID, _ := uuid.Parse(c.Param("userId"))

	if err := nc.NoteService.RevokeNoteShare(noteID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke share"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "access revoked"})
}
