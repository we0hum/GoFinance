package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func mustGetenv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("переменная %s не найдена в окружении", key)
	}
	return value
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Файл .env не найден")
	}

	connStr := mustGetenv("DATABASE_URL")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("База недоступна: %v", err)
	}

	fmt.Println("✅ Подключение к PostgreSQL установлено!")
}
