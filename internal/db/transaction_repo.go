package db

import (
	"GoFinance/internal/models"
	"context"
	"errors"
	"fmt"
	"strings"

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

	defer tx.Rollback()

	insertQuery := `
        INSERT INTO transactions (category_id, amount, note)
        VALUES ($1, $2, $3)
        RETURNING id, category_id, amount, note, created_at;
    `

	var tr models.Transaction
	if err := tx.GetContext(ctx, &tr, insertQuery, t.CategoryID, t.Amount, t.Note); err != nil {
		return models.Transaction{}, err
	}

	updateQuery := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
	`

	result, err := tx.ExecContext(ctx, updateQuery, t.Amount, accountID)

	if err != nil {
		return models.Transaction{}, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return models.Transaction{}, err
	}

	if rows == 0 {
		return models.Transaction{}, errors.New("account not found")
	}

	if err := tx.Commit(); err != nil {
		return models.Transaction{}, err
	}

	return tr, nil
}

func (r *TransactionRepo) List(ctx context.Context, filter TransactionFilter) ([]models.Transaction, error) {
	query := `
		SELECT id, category_id, amount, note, created_at
		FROM transactions	
	`

	args := []interface{}{}
	conditions := []string{}

	if filter.CategoryID != nil {
		args = append(args, *filter.CategoryID)

		conditions = append(conditions, fmt.Sprintf("category_id = $%d", len(args)))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	args = append(args, filter.Limit)

	query += fmt.Sprintf(
		" ORDER BY created_at DESC LIMIT $%d",
		len(args),
	)

	var txs []models.Transaction

	err := r.db.SelectContext(ctx, &txs, query, args...)

	if err != nil {
		return nil, err
	}

	return txs, err
}
