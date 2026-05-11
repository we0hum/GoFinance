package models

type Category struct {
	ID        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Type      string `db:"type" json:"type"`
	CreatedAt string `db:"created_at" json:"created_at"`
}
