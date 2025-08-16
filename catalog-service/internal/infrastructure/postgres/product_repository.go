package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/domain"
	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/infrastructure/postgres/dao"
)

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Save(ctx context.Context, p domain.Product) (int64, error) {
	query := `
		INSERT INTO products (name, description, price, category_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var id int64
	err := r.db.QueryRowContext(ctx, query, p.Name, p.Description, p.Price, p.CategoryID).Scan(&id)
	return id, err
}

func (r *productRepository) FindByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `SELECT id, name, description, price, category_id FROM products WHERE id = $1`

	var row dao.ProductRow
	err := r.db.GetContext(ctx, &row, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}

	product := mapToDomain(row)
	return &product, nil
}

func (r *productRepository) FindByIDs(ctx context.Context, ids []int64) ([]domain.Product, error) {
	const op = "productRepository.FindByIDs"

	if len(ids) == 0 {
		return []domain.Product{}, nil
	}

	query, args, err := sqlx.In("SELECT id, name, description, price, category_id FROM products WHERE id IN (?);", ids)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = r.db.Rebind(query)

	var products []domain.Product
	if err := r.db.SelectContext(ctx, &products, query, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	return products, nil
}

func (r *productRepository) FindAll(ctx context.Context) ([]domain.Product, error) {
	query := `SELECT id, name, description, price, category_id FROM products ORDER BY name ASC`

	var rows []dao.ProductRow
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}

	products := make([]domain.Product, 0, len(rows))
	for _, row := range rows {
		products = append(products, mapToDomain(row))
	}

	return products, nil
}

func (r *productRepository) FindByCategoryID(ctx context.Context, categoryID int64) ([]domain.Product, error) {
	query := `SELECT id, name, description, price, category_id FROM products WHERE category_id = $1 ORDER BY name ASC`

	var rows []dao.ProductRow
	if err := r.db.SelectContext(ctx, &rows, query, categoryID); err != nil {
		return nil, err
	}

	products := make([]domain.Product, 0, len(rows))
	for _, row := range rows {
		products = append(products, mapToDomain(row))
	}

	return products, nil
}

func (r *productRepository) Update(ctx context.Context, p domain.Product) error {
	query := `UPDATE products SET name = $1, description = $2, price = $3, category_id = $4 WHERE id = $5`
	res, err := r.db.ExecContext(ctx, query, p.Name, p.Description, p.Price, p.CategoryID, p.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func mapToDomain(p dao.ProductRow) domain.Product {
	return domain.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CategoryID:  p.CategoryID,
	}
}
