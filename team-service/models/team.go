package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name string    `gorm:"unique;not null" json:"name"`

	Members  []TeamMember  `json:"members" gorm:"foreignKey:TeamID"`
	Managers []TeamManager `json:"managers" gorm:"foreignKey:TeamID"`
}

func (t *Team) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}
