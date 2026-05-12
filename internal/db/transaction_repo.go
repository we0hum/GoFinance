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

func (r *TransactionRepo) Create(ctx context.Context, t models.Transaction) (models.Transaction, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return models.Transaction{}, err
	}

	query := `
        INSERT INTO transactions (category_id, amount, note)
        VALUES ($1, $2, $3)
        RETURNING id, category_id, amount, note, created_at;
    `

	var tr models.Transaction
	if err := r.db.GetContext(ctx, &tr, query, t.CategoryID, t.Amount, t.Note); err != nil {
		tx.Rollback()
		return models.Transaction{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Transaction{}, err
	}

	return tr, nil
}

func (r *TransactionRepo) CreateWithAccountUpdate(ctx context.Context, t models.Transaction, accountID int) (models.Transaction, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return models.Transaction{}, err
	}

	insertQuery := `
        INSERT INTO transactions (category_id, amount, note)
        VALUES ($1, $2, $3)
        RETURNING id, category_id, amount, note, created_at;
    `

	var tr models.Transaction
	if err := r.db.GetContext(ctx, &tr, insertQuery, t.CategoryID, t.Amount, t.Note); err != nil {
		tx.Rollback()
		return models.Transaction{}, err
	}

	updateQuery := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
	`

	if _, err := tx.ExecContext(ctx, updateQuery, t.Amount, accountID); err != nil {
		tx.Rollback()
		return models.Transaction{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Transaction{}, err
	}

	return tr, nil
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

func (r *TransactionRepo) ListByCategory(ctx context.Context, t models.Transaction) ([]models.Transaction, error) {
	var tx []models.Transaction

	err := r.db.SelectContext(ctx, &tx, `
        SELECT id, category_id, amount, note, created_at
        FROM transactions
        WHERE category_id = $1
        ORDER BY created_at DESC;
    `, t.CategoryID)

	return tx, err
}
