package services

import (
	"asset-service/models"
	"asset-service/repositories"

	"github.com/google/uuid"
)

type NoteService struct {
	NoteRepo *repositories.NoteRepository
}

// Khởi tạo NoteService
func NewNoteService(noteRepo *repositories.NoteRepository) *NoteService {
	return &NoteService{
		NoteRepo: noteRepo,
	}
}

// Tạo ghi chú mới
func (s *NoteService) CreateNote(note *models.Note) error {
	return s.NoteRepo.Create(note)
}

// Lấy ghi chú theo ID (bao gồm share info)
func (s *NoteService) GetNoteByID(noteID uuid.UUID) (*models.Note, error) {
	return s.NoteRepo.GetByID(noteID)
}

// Chia sẻ ghi chú cho người dùng khác
func (s *NoteService) ShareNoteWithUser(noteID, userID uuid.UUID, access string) error {
	return s.NoteRepo.ShareNote(noteID, userID, access)
}

// Thu hồi chia sẻ ghi chú với người dùng
func (s *NoteService) RevokeNoteShare(noteID, userID uuid.UUID) error {
	return s.NoteRepo.RevokeShare(noteID, userID)
}

// Kiểm tra quyền truy cập ghi chú
func (s *NoteService) HasAccess(note *models.Note, userID uuid.UUID) (string, bool) {
	if note.OwnerID == userID {
		return "write", true
	}

	access, err := s.NoteRepo.GetNoteAccess(note.ID, userID)
	if err != nil {
		return "", false
	}

	return access, true
}
