package models

type Report struct {
	TotalRevenue        int                `json:"total_revenue"`
	TotalTransactions   int                `json:"total_transactions"`
	BestSellingProducts BestSellingProduct `json:"best_selling_products"`
}

type BestSellingProduct struct {
	Name    string `json:"name"`
	QtySold int    `json:"qty_sold"`
}
