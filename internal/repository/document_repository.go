package repository

import (
	"github.com/arazumut/Lexa/internal/domain"
	"gorm.io/gorm"
)

type documentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) domain.DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) Create(doc *domain.Document) error {
	return r.db.Create(doc).Error
}

func (r *documentRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Document{}, id).Error
}

func (r *documentRepository) FindByID(id uint) (*domain.Document, error) {
	var doc domain.Document
	err := r.db.First(&doc, id).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *documentRepository) FindByCaseID(caseID uint) ([]domain.Document, error) {
	var docs []domain.Document
	// Yükleyen kullanıcı bilgisiyle beraber çekelim
	err := r.db.Preload("Uploader").Where("case_id = ?", caseID).Order("created_at desc").Find(&docs).Error
	if err != nil {
		return nil, err
	}
	return docs, nil
}
