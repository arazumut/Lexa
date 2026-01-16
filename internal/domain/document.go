package domain

import (
	"time"
)

// Document, dava veya iş dosyalarına eklenen evrakları temsil eder.
type Document struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CaseID      uint      `gorm:"index" json:"case_id"` // Hangi davaya ait?
	Case        Case      `json:"-"`                    // Relation
	UploaderID  uint      `json:"uploader_id"`          // Kim yükledi?
	Uploader    User      `gorm:"foreignKey:UploaderID" json:"uploader"`
	FileName    string    `json:"file_name"`            // Orijinal dosya adı (dilekce.pdf)
	FilePath    string    `json:"file_path"`            // Diskteki yolu (uploads/uuid.pdf)
	FileType    string    `json:"file_type"`            // Mime type (application/pdf)
	Category    string    `json:"category"`             // Dilekçe, Karar, Bilirkişi Raporu vb.
	Description string    `json:"description"`          // Ek açıklama
	CreatedAt   time.Time `json:"created_at"`
}

// DocumentRepository arayüzü
type DocumentRepository interface {
	Create(doc *Document) error
	Delete(id uint) error
	FindByID(id uint) (*Document, error)
	FindByCaseID(caseID uint) ([]Document, error)
}

// DocumentService arayüzü
type DocumentService interface {
	Upload(fileHeader interface{}, caseID uint, uploaderID uint, category, description string) (*Document, error) // fileHeader: multipart.FileHeader olacak ama interface tutuyoruz bağımlılık olmasın diye
	Delete(id uint) error
	GetListByCase(caseID uint) ([]Document, error)
	GetDocument(id uint) (*Document, error)
}
