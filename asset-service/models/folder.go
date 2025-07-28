package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Folder struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name    string    `gorm:"not null" json:"name"`
	OwnerID uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`

	Notes  []Note        `gorm:"foreignKey:FolderID" json:"notes,omitempty"`
	Shares []FolderShare `gorm:"foreignKey:FolderID" json:"shares,omitempty"`
}

func (u *Folder) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
