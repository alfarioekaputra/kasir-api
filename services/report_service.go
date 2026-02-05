package services

import (
	"labkoding.my.id/kasir-api/models"
	"labkoding.my.id/kasir-api/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{
		repo: repo,
	}
}

func (s *ReportService) TodayReport() (models.Report, error) {
	return s.repo.TodayReport()
}

func (s *ReportService) RangeReport(startDate, endDate string) (models.Report, error) {
	return s.repo.Range(startDate, endDate)
}
