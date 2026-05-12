-- +goose Up
CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions (category_id);

-- +goose Down
DROP INDEX IF EXISTS idx_transactions_category;