package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Kategori makanan"},
	{ID: 2, Name: "Minuman", Description: "Kategori minuman"},
	{ID: 3, Name: "Elektronik", Description: "Kategori elektronik"},
}

func GetAll(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(categories)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func FindCategoryById(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		slog.Error(err.Error())
	}

	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return nil
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
	return nil
}

func CreateCategory(w http.ResponseWriter, r *http.Request) error {
	var newCategory Category

	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		slog.Error(err.Error())
		return nil
	}

	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)

	return nil
}

func UpdateCategoryById(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		slog.Error(err.Error())
	}

	var updateCategory Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		slog.Error(err.Error())
		return nil
	}

	for i := range categories {
		if categories[i].ID == id {
			updateCategory.ID = id
			categories[i] = updateCategory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategory)
			return nil
		}
	}

	return nil
}

func DeleteCategoryById(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		slog.Error(err.Error())
	}

	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "success delete",
			})
			return nil
		}
	}

	http.Error(w, "category belum ada", http.StatusNotFound)
	return nil
}
