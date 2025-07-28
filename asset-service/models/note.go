package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Title    string    `gorm:"not null" json:"title"`
	Body     string    `gorm:"type:text" json:"body"`
	FolderID uuid.UUID `gorm:"type:uuid;not null" json:"folder_id"`
	OwnerID  uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`

	Shares []NoteShare `gorm:"foreignKey:NoteID" json:"shares,omitempty"`
}

func (u *Note) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
