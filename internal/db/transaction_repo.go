package db

import (
	"GoFinance/internal/models"
	"context"

	"github.com/jmoiron/sqlx"
)

type TransactionRepo struct {
	db *sqlx.DB
}

func NewTransactionRepo(db *sqlx.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) Create(ctx context.Context, categoryId int, amount float64, note string) (models.Transaction, error) {
	var tx models.Transaction

	query := `
        INSERT INTO transactions (category_id, amount, note)
        VALUES ($1, $2, $3)
        RETURNING id, category_id, amount, note, created_at;
    `

	if err := r.db.GetContext(ctx, &tx, query, categoryId, amount, note); err != nil {
		return models.Transaction{}, err
	}

	return tx, nil
}

func (r *TransactionRepo) List(ctx context.Context) ([]models.Transaction, error) {
	var txs []models.Transaction

	err := r.db.SelectContext(ctx, &txs, `
        SELECT id, category_id, amount, note, created_at
        FROM transactions
        ORDER BY created_at DESC;
    `)

	return txs, err
}

func (r *TransactionRepo) ListByCategory(ctx context.Context, categoryId int) ([]models.Transaction, error) {
	var tx []models.Transaction

	err := r.db.SelectContext(ctx, &tx, `
        SELECT id, category_id, amount, note, created_at
        FROM transactions
        WHERE category_id = $1
        ORDER BY created_at DESC;
    `, categoryId)

	return tx, err
}
