package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UUIDModel provides UUID primary key with auto-generation
// Embed this in your models instead of manually defining ID
type UUIDModel struct {
	ID string `json:"id" gorm:"type:char(36);primaryKey"`
}

// BeforeCreate hook to generate UUID before inserting
func (u *UUIDModel) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
