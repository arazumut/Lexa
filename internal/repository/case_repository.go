package repository

import (
	"github.com/arazumut/Lexa/internal/domain"
	"gorm.io/gorm"
)

type caseRepository struct {
	db *gorm.DB
}

// NewCaseRepository creates a new instance of CaseRepository
func NewCaseRepository(db *gorm.DB) domain.CaseRepository {
	return &caseRepository{db: db}
}

func (r *caseRepository) Create(c *domain.Case) error {
	return r.db.Create(c).Error
}

func (r *caseRepository) Update(c *domain.Case) error {
	return r.db.Save(c).Error
}

func (r *caseRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Case{}, id).Error
}

func (r *caseRepository) FindByID(id uint) (*domain.Case, error) {
	var c domain.Case
	// Müvekkil bilgisini de ("Client") preload ile çekiyoruz.
	err := r.db.Preload("Client").First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *caseRepository) FindAll(page, pageSize int, search string, clientID uint) ([]domain.Case, int64, int64, error) {
	var cases []domain.Case
	var totalCount int64
	var filteredCount int64

	// Base Query
	query := r.db.Model(&domain.Case{}).Preload("Client") // Listede müvekkil adını göstermek için

	// Client Filtresi (Eğer spesifik bir müvekkilin davaları isteniyorsa)
	if clientID != 0 {
		query = query.Where("client_id = ?", clientID)
	}

	// 1. Toplam Kayıt Sayısı
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, 0, err
	}

	// 2. Arama Filtresi (Search)
	if search != "" {
		searchTerm := "%" + search + "%"
		// Dosya No, Konu, Mahkeme veya Müvekkil Adı üzerinde arama yap
		query = query.Joins("LEFT JOIN clients ON clients.id = cases.client_id").
			Where("cases.file_number LIKE ? OR cases.title LIKE ? OR cases.court LIKE ? OR clients.name LIKE ?", 
			searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// 3. Filtrelenmiş Kayıt Sayısı
	if err := query.Count(&filteredCount).Error; err != nil {
		return nil, 0, 0, err
	}

	// 4. Veriyi Çek (Pagination)
	offset := (page - 1) * pageSize
	err := query.Order("cases.created_at desc").Offset(offset).Limit(pageSize).Find(&cases).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return cases, totalCount, filteredCount, nil
}

func (r *caseRepository) RecentCases(limit int) ([]domain.Case, error) {
	var cases []domain.Case
	// Sadece aktif davaları, son güncellenme tarihine göre getir
	err := r.db.Preload("Client").
		Where("status = ?", domain.CaseStatusActive).
		Order("updated_at desc").
		Limit(limit).
		Find(&cases).Error
	return cases, err
}
