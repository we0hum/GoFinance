package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")

	db, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	dir := "migrations"

	if len(os.Args) < 2 {
		log.Fatalf("Укажи команду: up | down | status | redo")
	}
	cmd := os.Args[1]

	if err := goose.Run(cmd, db, dir); err != nil {
		log.Fatalf("Ошибка при выполнении миграции: %v", err)
	}

	log.Println("✅ Команда успешно выполнена:", cmd)
}
