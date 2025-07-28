package services

import (
	"asset-service/models"
	"asset-service/repositories"

	"github.com/google/uuid"
)

type FolderService struct {
	FolderRepo *repositories.FolderRepository
	NoteRepo   *repositories.NoteRepository
}

// Khởi tạo FolderService
func NewFolderService(folderRepo *repositories.FolderRepository, noteRepo *repositories.NoteRepository) *FolderService {
	return &FolderService{
		FolderRepo: folderRepo,
		NoteRepo:   noteRepo,
	}
}

// Tạo mới folder
func (s *FolderService) CreateFolder(folder *models.Folder) error {
	return s.FolderRepo.Create(folder)
}

// Lấy folder theo ID (bao gồm notes và shares)
func (s *FolderService) GetFolderByID(id uuid.UUID) (*models.Folder, error) {
	return s.FolderRepo.GetByID(id)
}

// Chia sẻ folder và tất cả notes bên trong cho user khác
func (s *FolderService) ShareFolderWithUser(folderID, userID uuid.UUID, access string) error {
	// Bước 1: chia sẻ folder
	if err := s.FolderRepo.ShareFolder(folderID, userID, access); err != nil {
		return err
	}

	// Bước 2: lấy folder và danh sách notes bên trong
	folder, err := s.FolderRepo.GetByID(folderID)
	if err != nil {
		return err
	}

	// Bước 3: chia sẻ tất cả notes thuộc folder đó
	for _, note := range folder.Notes {
		if err := s.NoteRepo.ShareNote(note.ID, userID, access); err != nil {
			return err
		}
	}

	return nil
}

// Thu hồi chia sẻ folder và toàn bộ notes bên trong
func (s *FolderService) RevokeFolderShare(folderID, userID uuid.UUID) error {
	// Bước 1: xóa chia sẻ folder
	if err := s.FolderRepo.RevokeShare(folderID, userID); err != nil {
		return err
	}

	// Bước 2: lấy folder và notes
	folder, err := s.FolderRepo.GetByID(folderID)
	if err != nil {
		return err
	}

	// Bước 3: thu hồi từng ghi chú
	for _, note := range folder.Notes {
		if err := s.NoteRepo.RevokeShare(note.ID, userID); err != nil {
			return err
		}
	}

	return nil
}

// Kiểm tra quyền truy cập (owner hoặc được chia sẻ)
func (s *FolderService) HasAccess(folder *models.Folder, userID uuid.UUID) (string, bool) {
	if folder.OwnerID == userID {
		return "write", true
	}

	access, err := s.FolderRepo.GetFolderAccess(folder.ID, userID)
	if err != nil {
		return "", false
	}

	return access, true
}
