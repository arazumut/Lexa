package service

import (
	"errors"
	"strings"
	"time"

	"github.com/arazumut/Lexa/internal/domain"
)

type hearingService struct {
	repo     domain.HearingRepository
	caseRepo domain.CaseRepository
}

func NewHearingService(repo domain.HearingRepository, caseRepo domain.CaseRepository) domain.HearingService {
	return &hearingService{
		repo:     repo,
		caseRepo: caseRepo,
	}
}

func (s *hearingService) CreateHearing(h *domain.Hearing) error {
	// 1. Validasyonlar
	if h.CaseID == 0 {
		return errors.New("bir dava seçilmelidir")
	}
	if strings.TrimSpace(h.Title) == "" {
		return errors.New("duruşma konusu girilmelidir")
	}
	if h.Date.IsZero() {
		return errors.New("tarih ve saat seçilmelidir")
	}

	// Davanın varlığını kontrol et ve Lokasyonu otomatik doldur (Eğer boşsa)
	c, err := s.caseRepo.FindByID(h.CaseID)
	if err != nil {
		return errors.New("seçilen dava bulunamadı")
	}
	
	// Eğer duruşma yeri girilmemişse, davanın mahkemesini varsayılan yap
	if strings.TrimSpace(h.Location) == "" {
		h.Location = c.Court
	}

	// Varsayılan Statü
	if h.Status == "" {
		h.Status = domain.HearingStatusScheduled
	}

	return s.repo.Create(h)
}

func (s *hearingService) UpdateHearing(h *domain.Hearing) error {
	if strings.TrimSpace(h.Title) == "" {
		return errors.New("duruşma konusu boş olamaz")
	}
	return s.repo.Update(h)
}

func (s *hearingService) DeleteHearing(id uint) error {
	return s.repo.Delete(id)
}

func (s *hearingService) GetHearing(id uint) (*domain.Hearing, error) {
	return s.repo.FindByID(id)
}

func (s *hearingService) ListHearings(page, pageSize int, caseID uint) ([]domain.Hearing, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.FindAll(page, pageSize, caseID)
}

func (s *hearingService) GetUpcomingHearings(limit int) ([]domain.Hearing, error) {
	if limit < 1 {
		limit = 5
	}
	return s.repo.GetUpcoming(limit)
}
