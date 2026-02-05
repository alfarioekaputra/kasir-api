package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"labkoding.my.id/kasir-api/models"
	"labkoding.my.id/kasir-api/services"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

func (h *CategoryHandler) GetAllCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")
	categories, err := h.service.GetAllCategories(name)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	json.NewEncoder(w).Encode(categories)

}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var category models.CategoryRequest
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "ada kesalahan saat mengambil data", http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	err = h.service.CreateCategory(&category)
	if err != nil {
		http.Error(w, "ada kesalahan saat membuat kategori", http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	json.NewEncoder(w).Encode(category)

}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id tidak boleh kosong", http.StatusBadRequest)
		return
	}

	var category models.CategoryRequest
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "ada kesalahan saat mengambil data", http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	category.ID = id
	err = h.service.UpdateCategory(&category)
	if err != nil {
		http.Error(w, "ada kesalahan saat mengupdate category", http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	json.NewEncoder(w).Encode(category)

}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	err := h.service.DeleteCategory(id)
	if err != nil {
		http.Error(w, "ada kesalahan saat menghapus category", http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	json.NewEncoder(w).Encode("category berhasil dihapus")

}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		http.Error(w, "category tidak ditemukan", http.StatusNotFound)
		slog.Error(err.Error())
		return
	}

	json.NewEncoder(w).Encode(category)
}
