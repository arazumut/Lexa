package service

import (
	"errors"
	"strings"

	"github.com/arazumut/Lexa/internal/domain"
)

type caseService struct {
	repo       domain.CaseRepository
	clientRepo domain.ClientRepository // Müvekkil kontrolü için gerekli olabilir
}

// NewCaseService creates a new instance of CaseService
func NewCaseService(repo domain.CaseRepository, clientRepo domain.ClientRepository) domain.CaseService {
	return &caseService{
		repo:       repo,
		clientRepo: clientRepo,
	}
}

func (s *caseService) CreateCase(c *domain.Case) error {
	// 1. Validasyonlar
	if c.ClientID == 0 {
		return errors.New("bir müvekkil seçilmelidir")
	}
	if strings.TrimSpace(c.Title) == "" {
		return errors.New("dava başlığı boş olamaz")
	}
	
	// Müvekkilin varlığını kontrol et
	if _, err := s.clientRepo.FindByID(c.ClientID); err != nil {
		return errors.New("seçilen müvekkil bulunamadı")
	}

	// Varsayılan Değerler
	if c.Status == "" {
		c.Status = domain.CaseStatusActive
	}
	if c.Type == "" {
		c.Type = domain.CaseTypeOther
	}

	return s.repo.Create(c)
}

func (s *caseService) UpdateCase(c *domain.Case) error {
	if strings.TrimSpace(c.Title) == "" {
		return errors.New("dava başlığı boş olamaz")
	}
	return s.repo.Update(c)
}

func (s *caseService) DeleteCase(id uint) error {
	return s.repo.Delete(id)
}

func (s *caseService) GetCase(id uint) (*domain.Case, error) {
	return s.repo.FindByID(id)
}

func (s *caseService) ListCases(page, pageSize int, search string, clientID uint) ([]domain.Case, int64, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}
	return s.repo.FindAll(page, pageSize, search, clientID)
}

func (s *caseService) GetDashboardSummary() (map[string]interface{}, error) {
	// İleride dashboard için özet istatistikler buraya eklenebilir.
	// Örn: Toplam aktif dava sayısı, bu ayki duruşma sayısı vb.
	cases, err := s.repo.RecentCases(5)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"recent_cases": cases,
	}, nil
}
