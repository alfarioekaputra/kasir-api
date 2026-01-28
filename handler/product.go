package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"labkoding.my.id/kasir-api/models"
	"labkoding.my.id/kasir-api/services"
)

type Producthandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *Producthandler {
	return &Producthandler{
		service: service,
	}
}

func (h *Producthandler) GetAllProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	products, err := h.service.GetAllProducts()
	fmt.Println(products)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	json.NewEncoder(w).Encode(products)

}

func (h *Producthandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "ada kesalahan saat mengambil data", http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	err = h.service.CreateProduct(&product)
	if err != nil {
		http.Error(w, "ada kesalahan saat membuat produk", http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	json.NewEncoder(w).Encode(product)

}

func (h *Producthandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	product, err := h.service.GetProductByID(id)
	if err != nil {
		http.Error(w, "produk tidak ditemukan", http.StatusNotFound)
		slog.Error(err.Error())
		return
	}

	json.NewEncoder(w).Encode(product)

}

func (h *Producthandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id tidak boleh kosong", http.StatusBadRequest)
		return
	}
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "ada kesalahan saat mengambil data", http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	product.ID = id
	err = h.service.UpdateProduct(&product)
	if err != nil {
		http.Error(w, "ada kesalahan saat mengupdate produk", http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	json.NewEncoder(w).Encode(product)

}

func (h *Producthandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	err := h.service.DeleteProduct(id)
	if err != nil {
		http.Error(w, "ada kesalahan saat menghapus produk", http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	json.NewEncoder(w).Encode("produk berhasil dihapus")

}
