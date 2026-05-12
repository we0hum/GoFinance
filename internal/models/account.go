package models

type Account struct {
	ID      int     `db:"id" json:"id"`
	Name    string  `db:"name" json:"name"`
	Balance float64 `db:"balance" json:"balance"`
}
