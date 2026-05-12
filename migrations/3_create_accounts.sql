-- +goose Up
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    balance NUMERIC(10,2) NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS accounts;