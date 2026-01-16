package service

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arazumut/Lexa/internal/domain"
	"github.com/google/uuid"
)

type documentService struct {
	repo       domain.DocumentRepository
	uploadPath string
}

func NewDocumentService(repo domain.DocumentRepository, uploadPath string) domain.DocumentService {
	return &documentService{
		repo:       repo,
		uploadPath: uploadPath,
	}
}

func (s *documentService) Upload(fileHeader interface{}, caseID uint, uploaderID uint, category, description string) (*domain.Document, error) {
	// Type Assertion: Interface olarak gelen fileHeader'ı multipart.FileHeader'a çevir
	fh, ok := fileHeader.(*multipart.FileHeader)
	if !ok {
		return nil, errors.New("geçersiz dosya formatı")
	}

	// Dosyayı Aç
	file, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Benzersiz dosya adı oluştur (UUID)
	ext := filepath.Ext(fh.Filename)
	newFileName := uuid.New().String() + ext
	fullPath := filepath.Join(s.uploadPath, newFileName)

	// Hedef dosyayı oluştur
	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// İçeriği kopyala
	if _, err = io.Copy(dst, file); err != nil {
		return nil, err
	}

	// DB Kaydı Oluştur
	doc := &domain.Document{
		CaseID:      caseID,
		UploaderID:  uploaderID,
		FileName:    fh.Filename,
		FilePath:    "/assets/uploads/" + newFileName, // Web'den erişilebilir yol (Static FS)
		FileType:    strings.TrimPrefix(ext, "."),
		Category:    category,
		Description: description,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(doc); err != nil {
		// DB kaydı başarısız olursa yüklenen dosyayı sil
		os.Remove(fullPath)
		return nil, err
	}

	return doc, nil
}

func (s *documentService) Delete(id uint) error {
	// Önce dokümanı bul
	doc, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// 1. Veritabanından Sil
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// 2. Fiziksel Dosyayı Sil
	// FilePath web yolu olduğu için ("/assets/uploads/...") fiziksel yola çevirmemiz lazım
	// Ancak basitlik adına dosya adını çekip uploadPath ile birleştirebiliriz.
	// doc.FilePath: /assets/uploads/uuid.pdf -> sadece dosya adını alalım
	fileName := filepath.Base(doc.FilePath)
	physicalPath := filepath.Join(s.uploadPath, fileName)

	// Hata olsa bile DB'den silindiği için çok dert etmeyebiliriz ama loglamak iyi olur.
	os.Remove(physicalPath)

	return nil
}

func (s *documentService) GetListByCase(caseID uint) ([]domain.Document, error) {
	return s.repo.FindByCaseID(caseID)
}

func (s *documentService) GetDocument(id uint) (*domain.Document, error) {
	return s.repo.FindByID(id)
}
