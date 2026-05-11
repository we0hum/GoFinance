package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"GoFinance/internal/db"

	"github.com/jmoiron/sqlx"
)

type CategoryHandlers struct {
	repo *db.CategoryRepo
}

func NewCategoryHandlers(dbase *sqlx.DB) *CategoryHandlers {
	return &CategoryHandlers{repo: db.NewCategoryRepo(dbase)}
}

func (h *CategoryHandlers) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var in struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	in.Name = strings.TrimSpace(in.Name)

	if in.Name == "" || (in.Type != "income" && in.Type != "expense") {
		http.Error(w, "Name required, type must be income|expense", http.StatusBadRequest)
		return
	}

	cat, err := h.repo.Create(r.Context(), in.Name, in.Type)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cat)
}

func (h *CategoryHandlers) ListCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cats, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cats)
}
