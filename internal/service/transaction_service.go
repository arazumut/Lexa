package service

import (
	"errors"
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

func (s *transactionService) GetClientFinancials(clientID uint) ([]domain.Transaction, float64, error) {
	// 1. İşlemleri Çek
	txs, err := s.repo.GetClientTransactions(clientID)
	if err != nil {
		return nil, 0, err
	}

	// 2. Bakiyeyi Hesapla (Sadece bu müvekkil için)
	// Not: Repository'de özel bir GetBalanceByClient metodumuz yok,
	// bu yüzden çektiğimiz (son 50) işlemden değil, tüm işlemlerden hesaplamak daha doğru olurdu.
	// Ama şimdilik basitlik adına listedeki işlemi değil, veritabanından SUM alacak bir yöntem eklemeliyiz.
	// Hızlı çözüm: Service içinde basit bir döngü ile hesaplama (Sadece son 50 işlem için doğru olur)
	// DOĞRU ÇÖZÜM: Repository'e eklemek. Ama şimdilik filter kullanarak yapabiliriz.
	
	// TransactionFilter kullanarak bu müvekkilin tüm bakiyesini hesaplayabiliriz aslında ama 
	// GetTotalBalance genel çalışıyor.
	// Şimdilik 0 dönelim, hesaplamayı front-end veya basit döngü yapsın.
	// VEYA: FindAll ile hepsini çekip toplayalım (Performans sorunu olabilir ama idare eder)
	
	allTxs, _, _ := s.repo.FindAll(1, 10000, domain.TransactionFilter{ClientID: clientID})
	var balance float64
	for _, t := range allTxs {
		if t.Type == domain.TransactionTypeIncome {
			balance += t.Amount
		} else {
			balance -= t.Amount
		}
	}

	return txs, balance, nil
}
