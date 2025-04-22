package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPrice    = errors.New("invalid price")
	ErrInvalidProduct  = errors.New("invalid product")
	ErrProductNotFound = errors.New("product not found")
)

type Money float64

type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       Money
	SKU         string
	CategoryID  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Product) UpdatePrice(newPrice Money) error {
	if newPrice < 0 {
		return ErrInvalidPrice
	}
	p.Price = newPrice
	p.UpdatedAt = time.Now()

	return nil
}
