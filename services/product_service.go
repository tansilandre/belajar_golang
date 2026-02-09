package services

import (
	"andre_kasir_api/models"
	"andre_kasir_api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(searchName string) ([]models.Product, error) {
	return s.repo.GetAll(searchName)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Create(product *models.Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) Update(product *models.Product) error {
	return s.repo.Update(product)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
