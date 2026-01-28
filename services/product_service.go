package services

import (
	"labkoding.my.id/kasir-api/models"
	"labkoding.my.id/kasir-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.repo.GetAllProducts()
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.repo.CreateProduct(product)
}

func (s *ProductService) GetProductByID(id string) (*models.Product, error) {
	return s.repo.GetProductByID(id)
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	return s.repo.UpdateProduct(product)
}

func (s *ProductService) DeleteProduct(id string) error {
	return s.repo.DeleteProduct(id)
}
