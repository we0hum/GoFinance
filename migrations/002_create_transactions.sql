CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    category_id INT REFERENCES categories(id) ON DELETE CASCADE,
    amount NUMERIC(10,2) NOT NULL,
    note TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
