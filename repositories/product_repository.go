package repositories

import (
	"errors"

	"gorm.io/gorm"
	"labkoding.my.id/kasir-api/models"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) GetAllProducts(name string) ([]models.Product, error) {
	var products []models.Product

	query := r.db.Table("products").
		Select("products.id, products.name, products.description, products.price, products.stock, products.picture_url, products.category_id, categories.name as category_name").
		Joins("LEFT JOIN categories ON products.category_id = categories.id")

	if name != "" {
		query = query.Where("products.name ILIKE ?", "%"+name+"%")
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetProductByID(id string) (*models.Product, error) {
	var product models.Product

	err := r.db.Table("products").
		Select("products.id, products.name, products.description, products.price, products.stock, products.picture_url, products.category_id, categories.name as category_name").
		Joins("LEFT JOIN categories ON products.category_id = categories.id").
		Where("products.id = ?", id).
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	if err := r.db.Create(product).Error; err != nil {
		return err
	}

	// Get category name
	var category models.Category
	if err := r.db.First(&category, "id = ?", product.CategoryID).Error; err != nil {
		return err
	}
	product.CategoryName = category.Name

	return nil
}

func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	result := r.db.Model(&models.Product{}).
		Where("id = ?", product.ID).
		Updates(map[string]interface{}{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
			"category_id": product.CategoryID,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("produk tidak ditemukan")
	}

	// Update picture_url only if provided
	if product.PictureURL != nil {
		r.db.Model(&models.Product{}).Where("id = ?", product.ID).Update("picture_url", product.PictureURL)
	}

	return nil
}

func (r *ProductRepository) DeleteProduct(id string) error {
	result := r.db.Delete(&models.Product{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}
