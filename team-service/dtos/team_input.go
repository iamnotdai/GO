package dtos

import "github.com/google/uuid"

// AddUserToTeamInput dùng để thêm thành viên hoặc manager
type AddUserToTeamInput struct {
	UserID uuid.UUID `json:"userId" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}
