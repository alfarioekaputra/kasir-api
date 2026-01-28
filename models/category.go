package models

type Category struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type CategoryRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	ProductCount int       `json:"product_count"`
	Products     []Product `json:"products"`
}
