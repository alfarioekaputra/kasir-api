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

func (r *ProductRepository) GetAllProducts(name string) ([]models.Product, error) {
	var products []models.Product

	args := []interface{}{}
	query := "SELECT products.id, products.name, products.description, products.price, products.stock, products.picture_url, categories.id as category_id, categories.name as category_name FROM products LEFT JOIN categories ON products.category_id = categories.id"
	if name != "" {
		query += " WHERE products.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}
	rows, err := r.db.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.PictureURL, &product.CategoryID, &product.CategoryName); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetProductByID(id string) (*models.Product, error) {
	var product models.Product

	row := r.db.QueryRow("SELECT products.id, products.name, products.description, products.price, products.stock, products.picture_url, categories.id as category_id, categories.name as category_name FROM products LEFT JOIN categories ON products.category_id = categories.id WHERE products.id = $1", id)

	if err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.PictureURL, &product.CategoryID, &product.CategoryName); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	err := r.db.QueryRow("INSERT INTO products (name, description, price, stock, category_id, picture_url) VALUES ($1, $2, $3, $4, $5, $6) returning id, category_id", product.Name, product.Description, product.Price, product.Stock, product.CategoryID, product.PictureURL).Scan(&product.ID, &product.CategoryID)

	if err != nil {
		return err
	}

	err = r.db.QueryRow("SELECT name FROM categories WHERE id = $1", product.CategoryID).Scan(&product.CategoryName)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	result, err := r.db.Exec("UPDATE products SET name = $1, description = $2, price = $3, stock = $4, category_id = $5, picture_url = COALESCE($6, picture_url) WHERE id = $7", product.Name, product.Description, product.Price, product.Stock, product.CategoryID, product.PictureURL, product.ID)
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
	return nil
}

func (r *ProductRepository) DeleteProduct(id string) error {
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
	return nil
}
