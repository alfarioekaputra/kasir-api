package services

import (
	"labkoding.my.id/kasir-api/models"
	"labkoding.my.id/kasir-api/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) GetAllCategories() ([]models.CategoryResponse, error) {
	return s.repo.GetAllCategories()
}

func (s *CategoryService) CreateCategory(category *models.CategoryRequest) error {
	return s.repo.CreateCategory(category)
}

func (s *CategoryService) GetCategoryByID(id string) (*models.Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *CategoryService) UpdateCategory(category *models.CategoryRequest) error {
	return s.repo.UpdateCategory(category)
}

func (s *CategoryService) DeleteCategory(id string) error {
	return s.repo.DeleteCategory(id)
}
