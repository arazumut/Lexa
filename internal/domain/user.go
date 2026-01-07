package domain

import (
	"time"

	"gorm.io/gorm"
)

// User, sistemdeki kullanıcıları temsil eder (GORM Uyumlu).
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"` // Şifre hash'i
	Name      string         `gorm:"not null" json:"name"`
	Role      string         `gorm:"default:staff" json:"role"` // 'admin', 'staff'
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft Delete (Silindiğinde veri kaybolmaz)
}

// UserRepository, veritabanı işlemlerini soyutlar.
type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	GetByID(id uint) (*User, error)
}

// UserService, iş mantığını soyutlar.
type UserService interface {
	Register(email, password, name string) error
	Login(email, password string) (string, error)
}
