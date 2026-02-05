package repositories

import (
	"database/sql"

	"labkoding.my.id/kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{
		db: db,
	}
}

func (r *ReportRepository) TodayReport() (models.Report, error) {
	var report models.Report

	err := r.db.QueryRow("SELECT COALESCE(SUM(total_amount),0), COALESCE(COUNT(*),0) as total_transaction FROM transactions WHERE DATE(created_at) = CURRENT_DATE").Scan(&report.TotalRevenue, &report.TotalTransactions)
	if err != nil {
		return models.Report{}, err
	}

	rows, err := r.db.Query("SELECT p.name, COALESCE(SUM(td.quantity),0) FROM transaction_details td JOIN products p ON td.product_id = p.id WHERE td.transaction_id IN (SELECT id FROM transactions WHERE DATE(created_at) = CURRENT_DATE) GROUP BY p.name ORDER BY SUM(td.quantity) DESC LIMIT 1")
	if err != nil {
		return models.Report{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var productName string
		var quantity int
		if err := rows.Scan(&productName, &quantity); err != nil {
			return models.Report{}, err
		}
		report.BestSellingProducts = models.BestSellingProduct{
			Name:    productName,
			QtySold: quantity,
		}
	}

	return report, nil
}

func (r *ReportRepository) Range(startDate, endDate string) (models.Report, error) {
	var report models.Report

	err := r.db.QueryRow("SELECT COALESCE(SUM(total_amount),0), COALESCE(COUNT(*),0) as total_transaction FROM transactions WHERE DATE(created_at) between $1 and $2", startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransactions)
	if err != nil {
		return models.Report{}, err
	}

	rows, err := r.db.Query("SELECT p.name, COALESCE(SUM(td.quantity),0) FROM transaction_details td JOIN products p ON td.product_id = p.id WHERE td.transaction_id IN (SELECT id FROM transactions WHERE DATE(created_at) between $1 and $2) GROUP BY p.name ORDER BY SUM(td.quantity) DESC LIMIT 1", startDate, endDate)
	if err != nil {
		return models.Report{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var productName string
		var quantity int
		if err := rows.Scan(&productName, &quantity); err != nil {
			return models.Report{}, err
		}
		report.BestSellingProducts = models.BestSellingProduct{
			Name:    productName,
			QtySold: quantity,
		}
	}

	return report, nil
}
