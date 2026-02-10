package repositories

import (
	"fmt"

	"gorm.io/gorm"
	"labkoding.my.id/kasir-api/models"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	// Start transaction
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Process each item
		for _, item := range items {
			var product models.Product
			if err := tx.First(&product, "id = ?", item.ProductID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return fmt.Errorf("product with id %s not found", item.ProductID)
				}
				return err
			}

			// Check stock
			if product.Stock < item.Quantity {
				return fmt.Errorf("insufficient stock for product %s", product.Name)
			}

			subTotal := product.Price * item.Quantity
			totalAmount += subTotal

			// Update stock
			if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return err
			}

			details = append(details, models.TransactionDetail{
				ProductID:   item.ProductID,
				ProductName: product.Name,
				Quantity:    item.Quantity,
				Subtotal:    subTotal,
			})
		}

		// Create transaction
		transaction := models.Transaction{
			TotalAmount: totalAmount,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// Create transaction details
		for i := range details {
			details[i].TransactionID = transaction.ID
			if err := tx.Create(&details[i]).Error; err != nil {
				return err
			}
		}

		transaction.Details = details
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Fetch the created transaction
	var transaction models.Transaction
	if err := r.db.Preload("Details").Order("created_at DESC").First(&transaction).Error; err != nil {
		return nil, err
	}

	// Set product names for details
	for i := range transaction.Details {
		var product models.Product
		r.db.First(&product, "id = ?", transaction.Details[i].ProductID)
		transaction.Details[i].ProductName = product.Name
	}

	return &transaction, nil
}
