package repositories

import (
	"asset-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FolderRepository struct {
	DB *gorm.DB
}

func NewFolderRepository(db *gorm.DB) *FolderRepository {
	return &FolderRepository{DB: db}
}

func (r *FolderRepository) Create(folder *models.Folder) error {
	return r.DB.Create(folder).Error
}

func (r *FolderRepository) GetByID(id uuid.UUID) (*models.Folder, error) {
	var folder models.Folder
	err := r.DB.Preload("Notes").Preload("Shares").First(&folder, "id = ?", id).Error
	return &folder, err
}

func (r *FolderRepository) ShareFolder(folderID, userID uuid.UUID, access string) error {
	share := models.FolderShare{
		FolderID: folderID,
		UserID:   userID,
		Access:   access,
	}
	return r.DB.Save(&share).Error
}

func (r *FolderRepository) RevokeShare(folderID, userID uuid.UUID) error {
	return r.DB.Delete(&models.FolderShare{}, "folder_id = ? AND user_id = ?", folderID, userID).Error
}

func (r *FolderRepository) GetFolderAccess(folderID, userID uuid.UUID) (string, error) {
	var share models.FolderShare
	err := r.DB.First(&share, "folder_id = ? AND user_id = ?", folderID, userID).Error
	if err != nil {
		return "", err
	}
	return share.Access, nil
}
