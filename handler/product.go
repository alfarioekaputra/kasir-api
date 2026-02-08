package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

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

	name := r.URL.Query().Get("name")
	products, err := h.service.GetAllProducts(name)
	fmt.Println(products)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	json.NewEncoder(w).Encode(products)

}

// parseProductFromForm parses product data from multipart form or JSON body
func (h *Producthandler) parseProductFromForm(w http.ResponseWriter, r *http.Request) (*models.Product, error) {
	var product models.Product

	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "multipart/form-data") {
		// Set max upload size to 1MB
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			return nil, fmt.Errorf("ukuran file maksimal 1MB: %w", err)
		}

		product.Name = r.FormValue("name")
		if desc := r.FormValue("description"); desc != "" {
			product.Description = &desc
		}
		if p := r.FormValue("price"); p != "" {
			if v, err := strconv.Atoi(p); err == nil {
				product.Price = v
			}
		}
		if s := r.FormValue("stock"); s != "" {
			if v, err := strconv.Atoi(s); err == nil {
				product.Stock = v
			}
		}
		product.CategoryID = r.FormValue("category_id")

		file, header, err := r.FormFile("picture_url")
		if err == nil {
			defer file.Close()

			// Validate file size
			if header.Size > 1<<20 {
				return nil, fmt.Errorf("ukuran file maksimal 1MB")
			}

			// delegate upload to service for better separation of concerns
			url, err := h.service.UploadProductImage(r.Context(), file, header.Filename, header.Header.Get("Content-Type"))
			if err != nil {
				return nil, fmt.Errorf("gagal mengupload gambar: %w", err)
			}
			product.PictureURL = &url
		}
	} else {
		// fallback to JSON body
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			return nil, fmt.Errorf("ada kesalahan saat mengambil data: %w", err)
		}
	}

	return &product, nil
}

func (h *Producthandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	product, err := h.parseProductFromForm(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	if err := h.service.CreateProduct(product); err != nil {
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

	product, err := h.parseProductFromForm(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	product.ID = id
	if err := h.service.UpdateProduct(product); err != nil {
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
