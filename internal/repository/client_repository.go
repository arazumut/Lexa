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

func (r *clientRepository) FindAll(page, pageSize int, search string) ([]domain.Client, int64, error) {
	var clients []domain.Client
	var totalCount int64

	// Base Query
	query := r.db.Model(&domain.Client{})

	// Arama Filtresi (Search)
	// İsimde veya Kimlik Numarasında arama yapar.
	if search != "" {
		// SQLite LIKE işlemi varsayılan olarak case-insensitive'dir ancak garantiye almak için LOWER kullanılabilir.
		// Ancak performans için direkt LIKE kullanıyoruz, GORM bunu halleder.
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ? OR identity LIKE ?", searchTerm, searchTerm)
	}

	// 1. Toplam Kayıt Sayısını Hesapla (Sayfalama için gerekli)
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// 2. Veriyi Çek (Pagination)
	offset := (page - 1) * pageSize
	
	// En son eklenenler en üstte görünsün.
	err := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&clients).Error
	if err != nil {
		return nil, 0, err
	}

	return clients, totalCount, nil
}
