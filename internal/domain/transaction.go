package domain

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string

type PaymentMethod string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"

	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodBank     PaymentMethod = "bank"
	PaymentMethodCredit   PaymentMethod = "credit"
	PaymentMethodCheck    PaymentMethod = "check"
)

type Transaction struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	ClientID      *uint           `gorm:"index" json:"client_id,omitempty"`
	Client        *Client         `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	CaseID        *uint           `gorm:"index" json:"case_id,omitempty"`
	Case          *Case           `gorm:"foreignKey:CaseID" json:"case,omitempty"`
	Type          TransactionType `gorm:"size:20;not null;index" json:"type"`
	Category      string          `gorm:"size:100" json:"category"`
	Amount        float64         `gorm:"type:decimal(15,2);not null" json:"amount"`
	PaymentMethod PaymentMethod   `gorm:"size:50;default:'bank'" json:"payment_method"`
	Date          time.Time       `gorm:"index;not null" json:"date"`
	Description   string          `gorm:"type:text" json:"description"`
	InvoiceNo     string          `gorm:"size:50" json:"invoice_no"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}

type TransactionRepository interface {
	Create(t *Transaction) error
	Update(t *Transaction) error
	Delete(id uint) error
	FindByID(id uint) (*Transaction, error)
	FindAll(page, pageSize int, filter TransactionFilter) ([]Transaction, int64, error)
	GetTotalBalance() (float64, error)
	GetMonthlyStats() ([]MonthlyStat, error)
	GetClientTransactions(clientID uint) ([]Transaction, error)
}

type TransactionService interface {
	CreateTransaction(t *Transaction) error
	UpdateTransaction(t *Transaction) error
	DeleteTransaction(id uint) error
	GetTransaction(id uint) (*Transaction, error)
	ListTransactions(page, pageSize int, filter TransactionFilter) ([]Transaction, int64, error)
	GetDashboardFinancials() (float64, []MonthlyStat, error)
	GetClientFinancials(clientID uint) ([]Transaction, float64, error)
}

type TransactionFilter struct {
	Search    string
	Type      TransactionType
	ClientID  uint
	CaseID    uint
	StartDate *time.Time
	EndDate   *time.Time
}

type MonthlyStat struct {
	Month       string  `json:"month"`
	TotalIncome float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
}
