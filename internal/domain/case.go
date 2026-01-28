package domain

import (
	"time"

	"gorm.io/gorm"
)

type CaseStatus string
type CaseType string

const (
	CaseStatusActive   CaseStatus = "active"   // Devam Ediyor
	CaseStatusDecision CaseStatus = "decision" // Karar Aşamasında
	CaseStatusAppeal   CaseStatus = "appeal"   // İstinaf / Temyiz
	CaseStatusClosed   CaseStatus = "closed"   // Kapandı

	CaseTypeCriminal       CaseType = "criminal"        // Ceza Hukuku
	CaseTypeCivil          CaseType = "civil"           // Medeni Hukuk (Boşanma, Miras)
	CaseTypeCommercial     CaseType = "commercial"      // Ticaret Hukuku
	CaseTypeAdministrative CaseType = "administrative"  // İdare Hukuku
	CaseTypeLabor          CaseType = "labor"           // İş Hukuku
	CaseTypeEnforcement    CaseType = "enforcement"     // İcra Hukuku
	CaseTypeOther          CaseType = "other"           // Diğer
)

type Case struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ClientID    uint           `gorm:"not null;index" json:"client_id"`        // Hangi müvekkile ait?
	Client      Client         `gorm:"foreignKey:ClientID" json:"client"`      // Relationship
	Title       string         `gorm:"size:200;not null" json:"title"`         // Davanın kısa başlığı (Örn: X ile Y Boşanma Davası)
	FileNumber  string         `gorm:"size:50;index" json:"file_number"`       // Dosya / Esas No (2024/123 E.)
	Court       string         `gorm:"size:150" json:"court"`                  // Mahkeme Adı
	Judge       string         `gorm:"size:100" json:"judge"`                  // Hakim Adı (Opsiyonel)
	Type        CaseType       `gorm:"size:50;default:'other'" json:"type"`    // Dava Türü
	Status      CaseStatus     `gorm:"size:20;default:'active'" json:"status"` // Durum
	Description string         `gorm:"type:text" json:"description"`           // Detaylı Açıklama
	StartDate   time.Time      `json:"start_date"`                             // Dava Açılış Tarihi
	Hearings    []Hearing      `gorm:"foreignKey:CaseID" json:"hearings,omitempty"`
	Documents   []Document     `gorm:"foreignKey:CaseID" json:"documents,omitempty"`    // YENİ: Evraklar
	Transactions []Transaction `gorm:"foreignKey:CaseID" json:"transactions,omitempty"` // YENİ: Muhasebe
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type CaseRepository interface {
	Create(c *Case) error
	Update(c *Case) error
	Delete(id uint) error
	FindByID(id uint) (*Case, error)
	FindAll(page, pageSize int, search string, clientID uint) ([]Case, int64, int64, error)
	RecentCases(limit int) ([]Case, error)
	GetCaseStats() (map[string]int64, error)
}

// CaseService interface (Port)
type CaseService interface {
	CreateCase(c *Case) error
	UpdateCase(c *Case) error
	DeleteCase(id uint) error
	GetCase(id uint) (*Case, error)
	ListCases(page, pageSize int, search string, clientID uint) ([]Case, int64, int64, error)
	GetDashboardSummary() (map[string]interface{}, error)
	GetCaseStatistics() (map[string]int64, error)
}
