package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

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
	if err != nil {
		slog.Error(err.Error())
		return
	}

	json.NewEncoder(w).Encode(products)

}
