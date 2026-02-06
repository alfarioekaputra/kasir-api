package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"

	"labkoding.my.id/kasir-api/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) GetAllCategories(name string) ([]models.CategoryResponse, error) {
	var categories []models.CategoryResponse

	args := []interface{}{}
	// SQL efisien: 1 query dengan json_agg untuk mengumpulkan products per kategori
	// NOTE: Pastikan database yang digunakan adalah PostgreSQL
	selectPart := `
			SELECT
				c.id,
				c.name,
				c.description,
				COUNT(p.id) AS product_count,
				COALESCE(
					json_agg(
						json_build_object(
							'id', p.id,
							'name', p.name,
							'description', p.description,
							'price', p.price,
							'stock', p.stock,
							'category_id', p.category_id,
							'category_name', c.name
						)
					) FILTER (WHERE p.id IS NOT NULL),
					'[]'
				) AS products
			FROM categories c
			LEFT JOIN products p ON c.id = p.category_id
			`

	groupOrderPart := `
		GROUP BY c.id, c.name, c.description
		ORDER BY c.id;
		`

	query := selectPart

	if name != "" {
		query += " WHERE c.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	query += groupOrderPart

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.CategoryResponse
		var productsJSON []byte

		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.ProductCount, &productsJSON); err != nil {
			return nil, err
		}

		// Unmarshal produk JSON ke []models.Product
		if err := json.Unmarshal(productsJSON, &category.Products); err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoryByID(id string) (*models.Category, error) {
	var category models.Category

	row := r.db.QueryRow("SELECT name, description FROM categories WHERE id = $1", id)

	if err := row.Scan(&category.Name, &category.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category tidak ditemukan")
		}
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) CreateCategory(category *models.CategoryRequest) error {
	err := r.db.QueryRow("INSERT INTO categories (name, description) VALUES ($1, $2) returning id", category.Name, category.Description).Scan(&category.ID)
	return err
}

func (r *CategoryRepository) UpdateCategory(category *models.CategoryRequest) error {
	result, err := r.db.Exec("UPDATE categories SET name = $1, description = $2 WHERE id = $3", category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("category tidak ditemukan")
	}
	return err
}

func (r *CategoryRepository) DeleteCategory(id string) error {
	result, err := r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("category tidak ditemukan")
	}
	return err
}
