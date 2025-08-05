package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/domain"
	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/infrastructure/postgres/dao"
)

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Save(ctx context.Context, c domain.Category) (int64, error) {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id;`

	var id int64
	err := r.db.QueryRowContext(ctx, query, c.Name).Scan(&id)
	return id, err
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	query := `SELECT id, name FROM categories ORDER BY name ASC;`

	var rows []dao.CategoryRow
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}

	categories := make([]domain.Category, 0, len(rows))
	for _, row := range rows {
		categories = append(categories, domain.Category{
			ID:   row.ID,
			Name: row.Name,
		})
	}

	return categories, nil
}
