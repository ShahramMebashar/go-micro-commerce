package api

import (
	"microservice/services/product-service/internal/domain"
	"microservice/services/product-service/internal/infrastructure/validator"
	"time"

	"github.com/google/uuid"
)

type ProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
	CategoryID  string  `json:"category_id"`
}

// Validate validates the ProductRequest
func (p *ProductRequest) Validate(v *validator.Validator) uuid.UUID {
	// Validate required fields
	v.Required("name", p.Name)
	v.Required("description", p.Description)
	v.Required("sku", p.SKU)

	// Validate price
	v.MinValue("price", p.Price, 0.01)

	// Validate and parse category ID
	categoryID, _ := v.ValidUUID("category_id", p.CategoryID)

	return categoryID
}

// services/product-service/internal/infrastructure/api/dto.go

// ToModel converts a ProductRequest to a domain.Product
func (p *ProductRequest) ToModel() (*domain.Product, error) {
	categoryID, err := uuid.Parse(p.CategoryID)
	if err != nil {
		return nil, err
	}

	return &domain.Product{
		Name:        p.Name,
		Description: p.Description,
		Price:       domain.Money(p.Price),
		SKU:         p.SKU,
		CategoryID:  categoryID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
	CategoryID  string  `json:"category_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// FromModel converts a domain.Product to a ProductResponse
func ProductResponseFromModel(p *domain.Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       float64(p.Price),
		SKU:         p.SKU,
		CategoryID:  p.CategoryID.String(),
		CreatedAt:   p.CreatedAt.Format(time.RFC1123),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC1123),
	}
}

// Note: APIResponse has been moved to response.go
