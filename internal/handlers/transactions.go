package handlers

import (
	"GoFinance/internal/db"
	"GoFinance/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type TransactionHandlers struct {
	repo *db.TransactionRepo
}

func NewTransactionHandlers(dbase *sqlx.DB) *TransactionHandlers {
	return &TransactionHandlers{repo: db.NewTransactionRepo(dbase)}
}

func (h *TransactionHandlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var in struct {
		CategoryID int     `json:"category_id"`
		Amount     float64 `json:"amount"`
		Note       string  `json:"note"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tx, err := h.repo.Create(r.Context(), in.CategoryID, in.Amount, in.Note)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}

func (h *TransactionHandlers) ListTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	categoryIdStr := r.URL.Query().Get("category_id")

	var (
		categoryId int
		err        error
	)

	if categoryIdStr != "" {
		categoryId, err = strconv.Atoi(categoryIdStr)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}

	var (
		txs []models.Transaction
	)

	if categoryIdStr != "" {
		txs, err = h.repo.ListByCategory(r.Context(), categoryId)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		txs, err = h.repo.List(r.Context())
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txs)
}
