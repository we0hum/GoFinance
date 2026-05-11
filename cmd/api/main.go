package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"GoFinance/internal/db"
	"GoFinance/internal/handlers"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")

	dbase := db.MustConnect(dsn)

	defer dbase.Close()

	cats := handlers.NewCategoryHandlers(dbase)

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cats.ListCategories(w, r)
		case http.MethodPost:
			cats.CreateCategory(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	txs := handlers.NewTransactionHandlers(dbase)

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			txs.ListTransactions(w, r)
		case http.MethodPost:
			txs.CreateTransaction(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server running on http://localhost:%s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
