package models

type Transaction struct {
	ID         int     `db:"id" json:"id"`
	CategoryID int     `db:"category_id" json:"category_id"`
	Amount     float64 `db:"amount" json:"amount"`
	Note       string  `db:"note" json:"note"`
	CreatedAt  string  `db:"created_at" json:"created_at"`
}
