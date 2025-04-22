package interfaces

import (
	"context"
	"microservice/services/product-service/internal/domain"

	"github.com/google/uuid"
)

type Service interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, int, error)
	Create(ctx context.Context, product *domain.Product) error
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error

	Search(ctx context.Context, query string) ([]*domain.Product, error)
	GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.Product, error)
	CategoryExists(ctx context.Context, categoryID uuid.UUID) (bool, error)
}
