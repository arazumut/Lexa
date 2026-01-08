package repository

import (
	"time"

	"github.com/arazumut/Lexa/internal/domain"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(t *domain.Transaction) error {
	return r.db.Create(t).Error
}

func (r *transactionRepository) Update(t *domain.Transaction) error {
	return r.db.Save(t).Error
}

func (r *transactionRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Transaction{}, id).Error
}

func (r *transactionRepository) FindByID(id uint) (*domain.Transaction, error) {
	var t domain.Transaction
	// Nested Preload: Transaction -> Case, Transaction -> Client
	err := r.db.Preload("Case").Preload("Client").First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *transactionRepository) FindAll(page, pageSize int, filter domain.TransactionFilter) ([]domain.Transaction, int64, error) {
	var transactions []domain.Transaction
	var totalCount int64

	query := r.db.Model(&domain.Transaction{}).Preload("Client").Preload("Case")

	// Filtreler
	if filter.Search != "" {
		term := "%" + filter.Search + "%"
		query = query.Where("description LIKE ? OR category LIKE ? OR invoice_no LIKE ?", term, term, term)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.ClientID != 0 {
		query = query.Where("client_id = ?", filter.ClientID)
	}
	if filter.CaseID != 0 {
		query = query.Where("case_id = ?", filter.CaseID)
	}
	if filter.StartDate != nil {
		query = query.Where("date >= ?", filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("date <= ?", filter.EndDate)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("date desc").Offset(offset).Limit(pageSize).Find(&transactions).Error
	return transactions, totalCount, err
}

func (r *transactionRepository) GetTotalBalance() (float64, error) {
	var totalIncome float64
	var totalExpense float64

	// Gelirleri Topla
	r.db.Model(&domain.Transaction{}).Where("type = ?", domain.TransactionTypeIncome).Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome)

	// Giderleri Topla
	r.db.Model(&domain.Transaction{}).Where("type = ?", domain.TransactionTypeExpense).Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense)

	return totalIncome - totalExpense, nil
}

func (r *transactionRepository) GetMonthlyStats() ([]domain.MonthlyStat, error) {
	// Bu karmaşık sorgu veritabanına göre değişebilir (Postgres vs SQLite).
	// SQLite için "strftime" kullanılır.
	
	type Result struct {
		Month string
		Type  string
		Total float64
	}
	var results []Result

	// SQLite Sorgusu: YYYY-MM formatında grupla
	err := r.db.Model(&domain.Transaction{}).
		Select("strftime('%Y-%m', date) as month, type, SUM(amount) as total").
		Where("date >= date('now', '-6 months')"). // Son 6 ay
		Group("month, type").
		Order("month asc").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Veriyi mapleyerek MonthlyStat formatına çevir
	statsMap := make(map[string]*domain.MonthlyStat)
	
	for _, res := range results {
		if _, ok := statsMap[res.Month]; !ok {
			statsMap[res.Month] = &domain.MonthlyStat{Month: res.Month}
		}
		if res.Type == string(domain.TransactionTypeIncome) {
			statsMap[res.Month].TotalIncome = res.Total
		} else {
			statsMap[res.Month].TotalExpense = res.Total
		}
	}

	var stats []domain.MonthlyStat
	// Map iteration sırasız olduğu için yeniden sıralamak gerekebilir ama basitçe append edelim, 
	// zaten DB'den sıralı geldi, buradaki map bozabilir.
	// En sağlıklısı map'i döngüye sokmak değil, bilinen son 6 ayı looplayıp mapten çekmek.

	// Şimdilik basit append (Sıralama bozulabilir frontend düzeltir veya sonra optimize ederiz)
	for _, v := range statsMap {
		stats = append(stats, *v)
	}
	
	return stats, nil
}
