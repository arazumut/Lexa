package domain

import (
	"time"

	"gorm.io/gorm"
)

type HearingStatus string

const (
	HearingStatusScheduled HearingStatus = "scheduled"
	HearingStatusCompleted HearingStatus = "completed"
	HearingStatusPostponed HearingStatus = "postponed"
	HearingStatusCancelled HearingStatus = "cancelled"
)

type Hearing struct {
	ID          uint          `gorm:"primaryKey" json:"id"`
	CaseID      uint          `gorm:"not null;index" json:"case_id"`
	Case        Case          `gorm:"foreignKey:CaseID" json:"case"`
	Title       string        `gorm:"size:200;not null" json:"title"`
	Date        time.Time     `gorm:"index;not null" json:"date"`
	Location    string        `gorm:"size:200" json:"location"`
	Status      HearingStatus `gorm:"size:20;default:'scheduled'" json:"status"`
	Description string        `gorm:"type:text" json:"description"`
	Result      string        `gorm:"type:text" json:"result"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type HearingRepository interface {
	Create(h *Hearing) error
	Update(h *Hearing) error
	Delete(id uint) error
	FindByID(id uint) (*Hearing, error)
	FindAll(page, pageSize int, caseID uint) ([]Hearing, int64, error)
	GetUpcoming(limit int) ([]Hearing, error)
}
type HearingService interface {
	CreateHearing(h *Hearing) error
	UpdateHearing(h *Hearing) error
	DeleteHearing(id uint) error
	GetHearing(id uint) (*Hearing, error)
	ListHearings(page, pageSize int, caseID uint) ([]Hearing, int64, error)
	GetUpcomingHearings(limit int) ([]Hearing, error)
}
