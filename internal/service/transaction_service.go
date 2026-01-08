package service

import (
	"errors"
	"time"

	"github.com/arazumut/Lexa/internal/domain"
)

type transactionService struct {
	repo domain.TransactionRepository
}

func NewTransactionService(repo domain.TransactionRepository) domain.TransactionService {
	return &transactionService{
		repo: repo,
	}
}

func (s *transactionService) CreateTransaction(t *domain.Transaction) error {
	// 1. Validasyonlar
	if t.Amount <= 0 {
		return errors.New("tutar 0'dan büyük olmalıdır")
	}
	if t.Date.IsZero() {
		return errors.New("tarih seçilmelidir")
	}
	
	// Varsayılan
	if t.Type == "" {
		t.Type = domain.TransactionTypeIncome
	}
	if t.PaymentMethod == "" {
		t.PaymentMethod = domain.PaymentMethodBank
	}

	return s.repo.Create(t)
}

func (s *transactionService) UpdateTransaction(t *domain.Transaction) error {
	if t.Amount <= 0 {
		return errors.New("tutar 0'dan büyük olmalıdır")
	}
	return s.repo.Update(t)
}

func (s *transactionService) DeleteTransaction(id uint) error {
	return s.repo.Delete(id)
}

func (s *transactionService) GetTransaction(id uint) (*domain.Transaction, error) {
	return s.repo.FindByID(id)
}

func (s *transactionService) ListTransactions(page, pageSize int, filter domain.TransactionFilter) ([]domain.Transaction, int64, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 10 }
	
	return s.repo.FindAll(page, pageSize, filter)
}

func (s *transactionService) GetDashboardFinancials() (float64, []domain.MonthlyStat, error) {
	balance, err := s.repo.GetTotalBalance()
	if err != nil {
		return 0, nil, err
	}
	
	stats, err := s.repo.GetMonthlyStats()
	if err != nil {
		// Grafik verisi çekilemese bile bakiyeyi dön
		return balance, nil, nil
	}
	
	return balance, stats, nil
}
