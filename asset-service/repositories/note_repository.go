package repositories

import (
	"asset-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NoteRepository struct {
	DB *gorm.DB
}

func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{DB: db}
}

func (r *NoteRepository) Create(note *models.Note) error {
	return r.DB.Create(note).Error
}

func (r *NoteRepository) GetByID(id uuid.UUID) (*models.Note, error) {
	var note models.Note
	err := r.DB.Preload("Shares").First(&note, "id = ?", id).Error
	return &note, err
}

func (r *NoteRepository) ShareNote(noteID, userID uuid.UUID, access string) error {
	share := models.NoteShare{
		NoteID: noteID,
		UserID: userID,
		Access: access,
	}
	return r.DB.Save(&share).Error
}

func (r *NoteRepository) RevokeShare(noteID, userID uuid.UUID) error {
	return r.DB.Delete(&models.NoteShare{}, "note_id = ? AND user_id = ?", noteID, userID).Error
}

func (r *NoteRepository) GetNoteAccess(noteID, userID uuid.UUID) (string, error) {
	var share models.NoteShare
	err := r.DB.First(&share, "note_id = ? AND user_id = ?", noteID, userID).Error
	if err != nil {
		return "", err
	}
	return share.Access, nil
}
