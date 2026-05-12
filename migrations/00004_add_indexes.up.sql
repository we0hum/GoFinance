-- +goose Up
CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions (category_id);
