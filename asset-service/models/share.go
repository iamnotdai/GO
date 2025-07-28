package models

import (
	"github.com/google/uuid"
)

type FolderShare struct {
	FolderID uuid.UUID `gorm:"type:uuid;primaryKey" json:"folder_id"`
	UserID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	Access   string    `gorm:"type:text;not null" json:"access"` // "read" or "write"
}

type NoteShare struct {
	NoteID uuid.UUID `gorm:"type:uuid;primaryKey" json:"note_id"`
	UserID uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	Access string    `gorm:"type:text;not null" json:"access"` // "read" or "write"
}
