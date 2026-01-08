package repository

import (
	"github.com/arazumut/Lexa/internal/domain"
	"gorm.io/gorm"
)

type clientRepository struct {
	db *gorm.DB
}

// NewClientRepository, ClientRepository interface'ini implemente eden struct'ı döner.
func NewClientRepository(db *gorm.DB) domain.ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) Create(client *domain.Client) error {
	return r.db.Create(client).Error
}

func (r *clientRepository) Update(client *domain.Client) error {
	return r.db.Save(client).Error
}

func (r *clientRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Client{}, id).Error
}

func (r *clientRepository) FindByID(id uint) (*domain.Client, error) {
	var client domain.Client
	// First metodu, kayıt bulamazsa ErrRecordNotFound döner, bu bizim için beklenen bir durum.
	err := r.db.First(&client, id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *clientRepository) FindAll(page, pageSize int, search string) ([]domain.Client, int64, int64, error) {
	var clients []domain.Client
	var totalCount int64
	var filteredCount int64

	// 1. Önce veritabanındaki TOPLAM kayıt sayısını al (Filtresiz)
	// Bu, DataTables'ın "toplam X kayıt" diyebilmesi için şart.
	if err := r.db.Model(&domain.Client{}).Count(&totalCount).Error; err != nil {
		return nil, 0, 0, err
	}

	// Base Query
	query := r.db.Model(&domain.Client{})

	// 2. Arama Filtresi (Search)
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ? OR identity LIKE ?", searchTerm, searchTerm)
	}

	// 3. Filtrelenmiş Kayıt Sayısını Hesapla
	if err := query.Count(&filteredCount).Error; err != nil {
		return nil, 0, 0, err
	}

	// 4. Veriyi Çek (Pagination)
	offset := (page - 1) * pageSize
	
	// En son eklenenler en üstte görünsün.
	err := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&clients).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return clients, totalCount, filteredCount, nil
}

func (r *clientRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Client{}).Count(&count).Error
	return count, err
}

func (r *clientRepository) GetClientStats() (map[string]int64, error) {
	stats := make(map[string]int64)
	
	type Result struct {
		Type  string
		Count int64
	}
	var results []Result

	// GROUP BY type
	err := r.db.Model(&domain.Client{}).Select("type, count(*) as count").Group("type").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	for _, res := range results {
		stats[res.Type] = res.Count
	}
	return stats, nil
}

