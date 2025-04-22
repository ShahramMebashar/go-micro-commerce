package postgres

import (
	"context"
	"microservice/services/product-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresProductRepository struct {
	DB *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{
		DB: db,
	}
}

func (r *PostgresProductRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, int, error) {
	rows, err := r.DB.Query(ctx, "SELECT * FROM products LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.SKU,
			&product.CategoryID,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	var total int
	row := r.DB.QueryRow(ctx, "SELECT count(*) FROM  products")

	err = row.Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *PostgresProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product

	err := r.DB.QueryRow(ctx, "SELECT * FROM products WHERE id = $1", id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.SKU,
		&product.CategoryID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	defer r.DB.Close()

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *PostgresProductRepository) Create(ctx context.Context, product *domain.Product) error {
	return nil
}

func (r *PostgresProductRepository) Update(ctx context.Context, product *domain.Product) error {
	return nil
}

func (r *PostgresProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *PostgresProductRepository) Search(ctx context.Context, query string) ([]*domain.Product, error) {
	return nil, nil
}

func (r *PostgresProductRepository) GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.Product, error) {
	return nil, nil
}

func (r *PostgresProductRepository) CategoryExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool

	err := r.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", id).Scan(&exists)

	defer r.DB.Close()

	if err != nil {
		return false, err
	}

	return exists, nil
}
