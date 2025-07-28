package models

import "github.com/google/uuid"

type TeamManager struct {
	TeamID uuid.UUID `gorm:"type:uuid"`
	UserID uuid.UUID `gorm:"type:uuid"`
}
