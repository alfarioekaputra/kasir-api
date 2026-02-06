package services

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"labkoding.my.id/kasir-api/external"
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

func (s *ProductService) GetAllProducts(name string) ([]models.Product, error) {
	return s.repo.GetAllProducts(name)
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

// UploadProductImage uploads an image reader to configured R2 and returns the public URL.
// filename is used to preserve extension when generating the storage key.
func (s *ProductService) UploadProductImage(ctx context.Context, r io.Reader, filename, contentType string) (string, error) {
	ext := filepath.Ext(filename)
	key := fmt.Sprintf("products/%d%s", time.Now().UnixNano(), ext)
	url, err := external.UploadObject(ctx, key, r, contentType)
	if err != nil {
		return "", err
	}
	return url, nil
}
