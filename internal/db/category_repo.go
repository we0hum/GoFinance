package db

import (
	"GoFinance/internal/models"
	"context"

	"github.com/jmoiron/sqlx"
)

type CategoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(ctx context.Context, name, ctype string) (models.Category, error) {
	var cat models.Category

	query := `
        INSERT INTO categories (name, type)
        VALUES ($1, $2)
        RETURNING id, name, type, created_at;
    `

	if err := r.db.GetContext(ctx, &cat, query, name, ctype); err != nil {
		return models.Category{}, err
	}

	return cat, nil
}

func (r *CategoryRepo) List(ctx context.Context) ([]models.Category, error) {
	var cats []models.Category

	err := r.db.SelectContext(ctx, &cats, `
        SELECT id, name, type, created_at
        FROM categories
        ORDER BY id;
    `)

	return cats, err
}
