package repositories

import (
	"errors"

	"gorm.io/gorm"
	"labkoding.my.id/kasir-api/models"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) GetAllCategories(name string) ([]models.CategoryResponse, error) {
	var categories []models.Category

	query := r.db.Model(&models.Category{})
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}

	// Build response with products
	var responses []models.CategoryResponse
	for _, category := range categories {
		var products []models.Product
		r.db.Where("category_id = ?", category.ID).Find(&products)

		// Set category name for each product
		for i := range products {
			products[i].CategoryName = category.Name
		}

		responses = append(responses, models.CategoryResponse{
			ID:           category.ID,
			Name:         category.Name,
			Description:  category.Description,
			ProductCount: len(products),
			Products:     products,
		})
	}

	return responses, nil
}

func (r *CategoryRepository) GetCategoryByID(id string) (*models.Category, error) {
	var category models.Category

	err := r.db.First(&category, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category tidak ditemukan")
		}
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) CreateCategory(category *models.CategoryRequest) error {
	newCategory := models.Category{
		Name:        category.Name,
		Description: &category.Description,
	}
	if err := r.db.Create(&newCategory).Error; err != nil {
		return err
	}
	category.ID = newCategory.ID
	return nil
}

func (r *CategoryRepository) UpdateCategory(category *models.CategoryRequest) error {
	result := r.db.Model(&models.Category{}).Where("id = ?", category.ID).Updates(map[string]interface{}{
		"name":        category.Name,
		"description": category.Description,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("category tidak ditemukan")
	}
	return nil
}

func (r *CategoryRepository) DeleteCategory(id string) error {
	result := r.db.Delete(&models.Category{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("category tidak ditemukan")
	}
	return nil
}
