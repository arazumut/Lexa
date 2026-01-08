package domain

import (
	"time"

	"gorm.io/gorm"
)

// ClientType müvekkil tipini belirtir.
type ClientType string

const (
	ClientTypeIndividual ClientType = "individual" // Bireysel
	ClientTypeCorporate  ClientType = "corporate"  // Kurumsal
)

// Client, sistemdeki müvekkilleri temsil eder.
// Clean Architecture: Entities
type Client struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Type      ClientType     `gorm:"type:varchar(20);default:'individual';not null" json:"type"` // Bireysel / Kurumsal
	Name      string         `gorm:"size:150;not null;index" json:"name"`                        // Ad Soyad veya Şirket Ünvanı
	Identity  string         `gorm:"size:20;uniqueIndex" json:"identity"`                        // TC Kimlik No veya Vergi No
	Email     string         `gorm:"size:100" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Address   string         `gorm:"type:text" json:"address"`
	Notes     string         `gorm:"type:text" json:"notes"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ClientRepository, veritabanı erişim katmanını soyutlar (Port).
type ClientRepository interface {
	Create(client *Client) error
	Update(client *Client) error
	Delete(id uint) error
	FindByID(id uint) (*Client, error)
	// FindAll fonksiyonu sayfalama (pagination) ve arama (search) desteği ile döner.
	// returns: results, totalCount, filteredCount, error
	FindAll(page, pageSize int, search string) ([]Client, int64, int64, error)
}

// ClientService, iş mantığını soyutlar (Port).
type ClientService interface {
	CreateClient(client *Client) error
	UpdateClient(client *Client) error
	DeleteClient(id uint) error
	GetClient(id uint) (*Client, error)
	ListClients(page, pageSize int, search string) ([]Client, int64, int64, error)
}
