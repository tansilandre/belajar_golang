package services

import (
	"andre_kasir_api/models"
	"andre_kasir_api/repositories"
	"time"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(req *models.CheckoutRequest) (*models.Transaction, error) {
	return s.repo.Checkout(req)
}

func (s *TransactionService) GetDailyReport() (*models.SalesReport, error) {
	return s.repo.GetDailyReport(time.Now())
}

func (s *TransactionService) GetReportByDateRange(startDate, endDate time.Time) (*models.SalesReport, error) {
	return s.repo.GetReportByDateRange(startDate, endDate)
}
