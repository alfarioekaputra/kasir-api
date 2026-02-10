package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID          string              `json:"id" gorm:"type:char(36);primaryKey"`
	TotalAmount int                 `json:"total_amount" gorm:"not null"`
	CreatedAt   time.Time           `json:"created_at" gorm:"autoCreateTime"`
	Details     []TransactionDetail `json:"details" gorm:"foreignKey:TransactionID;references:ID"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

func (Transaction) TableName() string {
	return "transactions"
}

type TransactionDetail struct {
	ID            string `json:"id" gorm:"type:char(36);primaryKey"`
	TransactionID string `json:"transaction_id" gorm:"type:char(36);not null"`
	ProductID     string `json:"product_id" gorm:"type:char(36);not null"`
	ProductName   string `json:"product_name" gorm:"-"`
	Quantity      int    `json:"quantity" gorm:"not null"`
	Subtotal      int    `json:"subtotal" gorm:"not null"`
}

func (td *TransactionDetail) BeforeCreate(tx *gorm.DB) error {
	if td.ID == "" {
		td.ID = uuid.New().String()
	}
	return nil
}

func (TransactionDetail) TableName() string {
	return "transaction_details"
}

type CheckoutItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}
