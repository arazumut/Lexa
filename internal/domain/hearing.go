package domain

import (
	"time"

	"gorm.io/gorm"
)

type HearingStatus string

const (
	HearingStatusScheduled HearingStatus = "scheduled" // Planlandı
	HearingStatusCompleted HearingStatus = "completed" // Tamamlandı
	HearingStatusPostponed HearingStatus = "postponed" // Ertelendi
	HearingStatusCancelled HearingStatus = "cancelled" // İptal Edildi
)

// Hearing, bir dava dosyasına ait duruşmayı veya randevuyu temsil eder.
// Clean Architecture: Entity
type Hearing struct {
	ID          uint          `gorm:"primaryKey" json:"id"`
	CaseID      uint          `gorm:"not null;index" json:"case_id"`   // Hangi dava?
	Case        Case          `gorm:"foreignKey:CaseID" json:"case"`   // Relationship
	Title       string        `gorm:"size:200;not null" json:"title"`  // Duruşma Konusu (Örn: Tanık Dinlenmesi)
	Date        time.Time     `gorm:"index;not null" json:"date"`      // Duruşma Tarihi ve Saati
	Location    string        `gorm:"size:200" json:"location"`        // Yer (Mahkeme salonu vb. Genelde davanın mahkemesidir ama değişebilir)
	Status      HearingStatus `gorm:"size:20;default:'scheduled'" json:"status"`
	Description string        `gorm:"type:text" json:"description"`    // Açıklama / Hazırlık Notları
	Result      string        `gorm:"type:text" json:"result"`         // Duruşma Sonucu / Ara Karar
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// HearingRepository interface (Port)
type HearingRepository interface {
	Create(h *Hearing) error
	Update(h *Hearing) error
	Delete(id uint) error
	FindByID(id uint) (*Hearing, error)
	// FindAll fonksiyonu tarih aralığına göre de filtreleme yapabilir.
	FindAll(page, pageSize int, caseID uint) ([]Hearing, int64, error)
	// GetUpcoming dashboard ve bildirimler için yaklaşan duruşmaları getirir.
	// days: Kaç gün sonrasına kadar? (Örn: 7 gün)
	GetUpcoming(limit int) ([]Hearing, error)
}

// HearingService interface (Port)
type HearingService interface {
	CreateHearing(h *Hearing) error
	UpdateHearing(h *Hearing) error
	DeleteHearing(id uint) error
	GetHearing(id uint) (*Hearing, error)
	ListHearings(page, pageSize int, caseID uint) ([]Hearing, int64, error)
	GetUpcomingHearings(limit int) ([]Hearing, error)
}
