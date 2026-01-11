package domain

import (
	"time"

	"gorm.io/gorm"
)

type ClientType string

const (
	ClientTypeIndividual ClientType = "individual"
	ClientTypeCorporate  ClientType = "corporate"
)

type Client struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Type      ClientType     `gorm:"type:varchar(20);default:'individual';not null" json:"type"`
	Name      string         `gorm:"size:150;not null;index" json:"name"`
	Identity  string         `gorm:"size:20;uniqueIndex" json:"identity"`
	Email     string         `gorm:"size:100" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Address   string         `gorm:"type:text" json:"address"`
	Notes     string         `gorm:"type:text" json:"notes"`
	Cases     []Case         `gorm:"foreignKey:ClientID" json:"cases,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ClientRepository interface {
	Create(client *Client) error
	Update(client *Client) error
	Delete(id uint) error
	FindByID(id uint) (*Client, error)
	FindAll(page, pageSize int, search string) ([]Client, int64, int64, error)
	Count() (int64, error)
	GetClientStats() (map[string]int64, error)
}

type ClientService interface {
	CreateClient(client *Client) error
	UpdateClient(client *Client) error
	DeleteClient(id uint) error
	GetClient(id uint) (*Client, error)
	ListClients(page, pageSize int, search string) ([]Client, int64, int64, error)
	GetTotalCount() (int64, error)
	GetClientStatistics() (map[string]int64, error)
}

