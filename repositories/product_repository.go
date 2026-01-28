package repositories

import (
	"database/sql"
	"errors"

	"labkoding.my.id/kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product

	rows, err := r.db.Query("SELECT id, name, description, price, stock FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetProductByID(id int) (*models.Product, error) {
	var product models.Product

	row := r.db.QueryRow("SELECT id, name, description, price, stock FROM products WHERE id = $1", id)

	if err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	_, err := r.db.Exec("INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4)", product.Name, product.Description, product.Price, product.Stock)
	return err
}

func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	result, err := r.db.Exec("UPDATE products SET name = $1, description = $2, price = $3, stock = $4 WHERE id = $5", product.Name, product.Description, product.Price, product.Stock, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}
	return err
}

func (r *ProductRepository) DeleteProduct(id int) error {
	result, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}
	return err
}
