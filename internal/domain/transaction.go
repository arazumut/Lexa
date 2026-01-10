package domain

import (
	"time"

	"gorm.io/gorm"
)

// TransactionType işlemin yönünü belirtir.
type TransactionType string

// PaymentMethod ödeme yöntemini belirtir.
type PaymentMethod string

const (
	TransactionTypeIncome  TransactionType = "income"  // Gelir (Tahsilat)
	TransactionTypeExpense TransactionType = "expense" // Gider (Masraf)

	PaymentMethodCash     PaymentMethod = "cash"     // Nakit
	PaymentMethodBank     PaymentMethod = "bank"     // Havale/EFT
	PaymentMethodCredit   PaymentMethod = "credit"   // Kredi Kartı
	PaymentMethodCheck    PaymentMethod = "check"    // Çek/Senet
)

// Transaction, sistemdeki tüm finansal hareketleri temsil eder.
// Clean Architecture: Entity
type Transaction struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	ClientID      *uint           `gorm:"index" json:"client_id,omitempty"`       // Opsiyonel (Genel gider olabilir)
	Client        *Client         `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	CaseID        *uint           `gorm:"index" json:"case_id,omitempty"`         // Opsiyonel (Davadan bağımsız olabilir)
	Case          *Case           `gorm:"foreignKey:CaseID" json:"case,omitempty"`
	Type          TransactionType `gorm:"size:20;not null;index" json:"type"`     // Gelir / Gider
	Category      string          `gorm:"size:100" json:"category"`               // Örn: Vekalet Ücreti, Harç, Kira, Maaş
	Amount        float64         `gorm:"type:decimal(15,2);not null" json:"amount"` // Tutar
	PaymentMethod PaymentMethod   `gorm:"size:50;default:'bank'" json:"payment_method"`
	Date          time.Time       `gorm:"index;not null" json:"date"`             // İşlem Tarihi
	Description   string          `gorm:"type:text" json:"description"`           // Açıklama
	InvoiceNo     string          `gorm:"size:50" json:"invoice_no"`              // Fatura/Makbuz No (Varsa)
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TransactionRepository interface (Port)
type TransactionRepository interface {
	Create(t *Transaction) error
	Update(t *Transaction) error
	Delete(id uint) error
	FindByID(id uint) (*Transaction, error)
	// FindAll gelişmiş filtreleme ile (Tarih aralığı, Tip, Müvekkil, Dava)
	FindAll(page, pageSize int, filter TransactionFilter) ([]Transaction, int64, error)
	// GetTotalBalance Toplam kasayı (Gelir - Gider) hesaplar
	GetTotalBalance() (float64, error)
	// GetMonthlyStats Son 6 ayın gelir/gider grafiği için veri döner
	GetMonthlyStats() ([]MonthlyStat, error)
	GetClientTransactions(clientID uint) ([]Transaction, error)
}

// TransactionService interface (Port)
type TransactionService interface {
	CreateTransaction(t *Transaction) error
	UpdateTransaction(t *Transaction) error
	DeleteTransaction(id uint) error
	GetTransaction(id uint) (*Transaction, error)
	ListTransactions(page, pageSize int, filter TransactionFilter) ([]Transaction, int64, error)
	GetDashboardFinancials() (float64, []MonthlyStat, error)
	GetClientFinancials(clientID uint) ([]Transaction, float64, error)
}

// DTOs & Helper Structs

type TransactionFilter struct {
	Search    string
	Type      TransactionType
	ClientID  uint
	CaseID    uint
	StartDate *time.Time
	EndDate   *time.Time
}

type MonthlyStat struct {
	Month       string  `json:"month"` // "2024-01"
	TotalIncome float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
}
