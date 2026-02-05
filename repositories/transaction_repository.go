package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"labkoding.my.id/kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	stmtProd, err := tx.Prepare("select name, price, stock from products where id = $1")
	if err != nil {
		return nil, err
	}
	defer stmtProd.Close()

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := stmtProd.QueryRow(item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %s not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subTotal := productPrice * item.Quantity
		totalAmount += subTotal

		_, err = tx.Exec("update products set stock = stock - $1 where id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subTotal,
		})
	}

	var transactionID string
	var createdAt time.Time
	err = tx.QueryRow("insert into transactions (total_amount) values ($1) returning id, created_at", totalAmount).Scan(&transactionID, &createdAt)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("insert into transaction_details (transaction_id, product_id, quantity, subtotal) values ($1, $2, $3, $4) returning id")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for i := range details {
		details[i].TransactionID = transactionID

		var detailID string
		err = stmt.QueryRow(transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).Scan(&detailID)
		if err != nil {
			return nil, err
		}
		details[i].ID = detailID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	fmt.Println(createdAt)
	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
		CreatedAt:   createdAt,
	}, nil
}
