package main

import (
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")

	db, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	goose.SetDialect("postgres")

	dir := "migrations"

	if len(os.Args) < 2 {
		log.Fatalf("Укажи команду: up | down | status | redo")
	}
	cmd := os.Args[1]

	switch cmd {
	case "up":
		err = goose.Up(db, dir)
	case "down":
		err = goose.Down(db, dir)
	case "status":
		err = goose.Status(db, dir)
	case "redo":
		err = goose.Redo(db, dir)
	default:
		log.Fatalf("Неизвестная команда: %s", cmd)
	}

	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	log.Println("✅ Команда успешно выполнена:", cmd)
}
