package service

import (
	"context"
	"product-service/internal/models"
	"product-service/internal/repository"
)

type ProductService interface {
	GetAll(ctx context.Context, limit, offset int, search string, status *int) ([]models.Product, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Product, error)
	GetByStatusActive(ctx context.Context, limit, offset int, search string) ([]models.Product, int64, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, id uint, product *models.Product) error
	Delete(ctx context.Context, id uint) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo}
}

func (s *productService) GetAll(ctx context.Context, limit, offset int, search string, status *int) ([]models.Product, int64, error) {
	return s.repo.GetAll(ctx, limit, offset, search, status)
}

func (s *productService) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) GetByStatusActive(ctx context.Context, limit, offset int, search string) ([]models.Product, int64, error) {
	return s.repo.GetByStatusActive(ctx, limit, offset, search)
}

func (s *productService) Create(ctx context.Context, product *models.Product) error {
	return s.repo.Create(ctx, product)
}

func (s *productService) Update(ctx context.Context, id uint, product *models.Product) error {
	product.ID = id
	return s.repo.Update(ctx, product)
}

func (s *productService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
