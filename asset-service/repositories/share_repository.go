package repositories

import (
	"asset-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShareRepository struct {
	DB *gorm.DB
}

func NewShareRepository(db *gorm.DB) *ShareRepository {
	return &ShareRepository{DB: db}
}

func (r *ShareRepository) GetAllSharedFolders(userID uuid.UUID) ([]models.Folder, error) {
	var folders []models.Folder
	err := r.DB.Joins("JOIN folder_shares ON folder_shares.folder_id = folders.id").
		Where("folder_shares.user_id = ?", userID).Find(&folders).Error
	return folders, err
}

func (r *ShareRepository) GetAllSharedNotes(userID uuid.UUID) ([]models.Note, error) {
	var notes []models.Note
	err := r.DB.Joins("JOIN note_shares ON note_shares.note_id = notes.id").
		Where("note_shares.user_id = ?", userID).Find(&notes).Error
	return notes, err
}
