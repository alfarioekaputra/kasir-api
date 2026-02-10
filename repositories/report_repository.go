package repositories

import (
	"gorm.io/gorm"
	"labkoding.my.id/kasir-api/models"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{
		db: db,
	}
}

func (r *ReportRepository) TodayReport() (models.Report, error) {
	var report models.Report

	// Get total revenue and transactions for today
	err := r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(total_amount), 0) as total_revenue, COALESCE(COUNT(*), 0) as total_transactions").
		Where("DATE(created_at) = CURRENT_DATE").
		Scan(&report).Error
	if err != nil {
		return models.Report{}, err
	}

	// Get best selling product for today
	var bestProduct struct {
		Name    string
		QtySold int
	}

	err = r.db.Table("transaction_details td").
		Select("p.name, COALESCE(SUM(td.quantity), 0) as qty_sold").
		Joins("JOIN products p ON td.product_id = p.id").
		Where("td.transaction_id IN (?)",
			r.db.Model(&models.Transaction{}).
				Select("id").
				Where("DATE(created_at) = CURRENT_DATE"),
		).
		Group("p.name").
		Order("qty_sold DESC").
		Limit(1).
		Scan(&bestProduct).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return models.Report{}, err
	}

	report.BestSellingProducts = models.BestSellingProduct{
		Name:    bestProduct.Name,
		QtySold: bestProduct.QtySold,
	}

	return report, nil
}

func (r *ReportRepository) Range(startDate, endDate string) (models.Report, error) {
	var report models.Report

	// Get total revenue and transactions for date range
	err := r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(total_amount), 0) as total_revenue, COALESCE(COUNT(*), 0) as total_transactions").
		Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate).
		Scan(&report).Error
	if err != nil {
		return models.Report{}, err
	}

	// Get best selling product for date range
	var bestProduct struct {
		Name    string
		QtySold int
	}

	err = r.db.Table("transaction_details td").
		Select("p.name, COALESCE(SUM(td.quantity), 0) as qty_sold").
		Joins("JOIN products p ON td.product_id = p.id").
		Where("td.transaction_id IN (?)",
			r.db.Model(&models.Transaction{}).
				Select("id").
				Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate),
		).
		Group("p.name").
		Order("qty_sold DESC").
		Limit(1).
		Scan(&bestProduct).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return models.Report{}, err
	}

	report.BestSellingProducts = models.BestSellingProduct{
		Name:    bestProduct.Name,
		QtySold: bestProduct.QtySold,
	}

	return report, nil
}
