package repository

import (
	"time"

	"github.com/arazumut/Lexa/internal/domain"
	"gorm.io/gorm"
)

type hearingRepository struct {
	db *gorm.DB
}

func NewHearingRepository(db *gorm.DB) domain.HearingRepository {
	return &hearingRepository{db: db}
}

func (r *hearingRepository) Create(h *domain.Hearing) error {
	return r.db.Create(h).Error
}

func (r *hearingRepository) Update(h *domain.Hearing) error {
	return r.db.Save(h).Error
}

func (r *hearingRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Hearing{}, id).Error
}

func (r *hearingRepository) FindByID(id uint) (*domain.Hearing, error) {
	var h domain.Hearing
	// Duruşmanın hangi davaya, o davanın da hangi müvekkile ait olduğunu bilmemiz gerekebilir.
	// Nested Preload: Hearing -> Case -> Client
	err := r.db.Preload("Case.Client").First(&h, id).Error
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *hearingRepository) FindAll(page, pageSize int, caseID uint) ([]domain.Hearing, int64, error) {
	var hearings []domain.Hearing
	var totalCount int64

	query := r.db.Model(&domain.Hearing{}).Preload("Case.Client")

	if caseID != 0 {
		query = query.Where("case_id = ?", caseID)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	// Tarihe göre sırala (En yakın tarih en üstte)
	err := query.Order("date asc").Offset(offset).Limit(pageSize).Find(&hearings).Error
	if err != nil {
		return nil, 0, err
	}

	return hearings, totalCount, nil
}

func (r *hearingRepository) GetUpcoming(limit int) ([]domain.Hearing, error) {
	var hearings []domain.Hearing
	now := time.Now()

	// Bugünden sonraki, planlanmış duruşmaları getir
	err := r.db.Preload("Case").
		Where("date >= ? AND status = ?", now, domain.HearingStatusScheduled).
		Order("date asc"). // En yakını en önce getir
		Limit(limit).
		Find(&hearings).Error

	return hearings, err
}
