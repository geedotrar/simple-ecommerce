package repository

import (
	"context"
	"fmt"
	"product-service/config"
	"product-service/internal/models"
)

type ProductRepository interface {
	GetAll(ctx context.Context, limit, offset int, search string, status *int) ([]models.Product, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Product, error)
	GetByStatusActive(ctx context.Context, limit, offset int, search string) ([]models.Product, int64, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uint) error
}

type productRepository struct {
	db config.GormPostgres
}

func NewProductRepository(db config.GormPostgres) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAll(ctx context.Context, limit, offset int, search string, status *int) ([]models.Product, int64, error) {
	conn := r.db.GetConnection()
	var products []models.Product
	var total int64

	query := conn.WithContext(ctx).Model(&models.Product{}).Where("deleted_at IS NULL")

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?  ", "%"+search+"%", "%"+search+"%")
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	conn := r.db.GetConnection()
	var product models.Product
	err := conn.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&product).Error
	return &product, err
}

func (r *productRepository) GetByStatusActive(ctx context.Context, limit, offset int, search string) ([]models.Product, int64, error) {
	conn := r.db.GetConnection()
	var products []models.Product
	var total int64

	query := conn.WithContext(ctx).Model(&models.Product{}).Where("status = ? AND deleted_at IS NULL", 1)

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?  ", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	conn := r.db.GetConnection()
	return conn.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	conn := r.db.GetConnection()
	return conn.WithContext(ctx).Create(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	conn := r.db.GetConnection()

	var product models.Product
	err := conn.WithContext(ctx).Unscoped().First(&product, id).Error
	if err != nil {
		return err
	}

	if product.DeletedAt.Valid {
		return fmt.Errorf("product already deleted")
	}

	return conn.WithContext(ctx).Delete(&product).Error
}
