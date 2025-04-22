package application

import (
	"context"
	"microservice/pkg/logger"
	"microservice/services/product-service/internal/domain"
	"microservice/services/product-service/internal/interfaces"

	"github.com/google/uuid"
)

type ProductService struct {
	repo   interfaces.ProductRepository
	logger logger.Logger
}

func NewProductService(repo interfaces.ProductRepository, logger logger.Logger) *ProductService {
	return &ProductService{
		repo:   repo,
		logger: logger,
	}
}

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) CategoryExists(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repo.CategoryExists(ctx, id)
}

func (s *ProductService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, int, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *ProductService) Create(ctx context.Context, product *domain.Product) error {
	if product.Name == "" {
		return domain.ErrInvalidProduct
	}

	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}

	return s.repo.Create(ctx, product)
}

func (s *ProductService) Update(ctx context.Context, product *domain.Product) error {
	return s.repo.Update(ctx, product)
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *ProductService) GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.Product, error) {
	return s.repo.GetByCategory(ctx, categoryID)
}

func (s *ProductService) Search(ctx context.Context, query string) ([]*domain.Product, error) {
	return s.repo.Search(ctx, query)
}
