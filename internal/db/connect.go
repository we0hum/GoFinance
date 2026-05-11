package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func MustConnect(dsn string) *sqlx.DB {
	db := sqlx.MustConnect("postgres", dsn)
	return db
}
