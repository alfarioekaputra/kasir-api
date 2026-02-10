package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID           string    `json:"id" gorm:"type:char(36);primaryKey"`
	Name         string    `json:"name" gorm:"type:varchar(255);not null"`
	Description  *string   `json:"description" gorm:"type:text"`
	Price        int       `json:"price" gorm:"not null"`
	Stock        int       `json:"stock" gorm:"not null"`
	CategoryID   string    `json:"category_id" gorm:"type:char(36);not null"`
	Category     *Category `json:"-" gorm:"foreignKey:CategoryID;references:ID"`
	CategoryName string    `json:"category_name" gorm:"-"`
	PictureURL   *string   `json:"picture_url,omitempty" gorm:"type:text"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

func (Product) TableName() string {
	return "products"
}
