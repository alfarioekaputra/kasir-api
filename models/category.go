package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID          string  `json:"id" gorm:"type:char(36);primaryKey"`
	Name        string  `json:"name" gorm:"type:varchar(255);not null;unique"`
	Description *string `json:"description" gorm:"type:text"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

func (Category) TableName() string {
	return "categories"
}

type CategoryRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	ProductCount int       `json:"product_count"`
	Products     []Product `json:"products"`
}
