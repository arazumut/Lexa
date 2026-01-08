package service

import (
	"errors"
	"strings"

	"github.com/arazumut/Lexa/internal/domain"
)

type clientService struct {
	repo domain.ClientRepository
}

// NewClientService, ClientService interface'ini implemente eden struct'ı döner.
// Dependency Injection: Repository servise enjekte edilir.
func NewClientService(repo domain.ClientRepository) domain.ClientService {
	return &clientService{
		repo: repo,
	}
}

func (s *clientService) CreateClient(client *domain.Client) error {
	// 1. Validasyonlar
	if strings.TrimSpace(client.Name) == "" {
		return errors.New("müvekkil adı boş olamaz")
	}

	// Email formatı vb. validasyonlar buraya eklenebilir.

	// 2. Varsayılan Değerler
	if client.Type == "" {
		client.Type = domain.ClientTypeIndividual
	}

	// 3. Veritabanı Kaydı
	return s.repo.Create(client)
}

func (s *clientService) UpdateClient(client *domain.Client) error {
	if strings.TrimSpace(client.Name) == "" {
		return errors.New("müvekkil adı boş olamaz")
	}
	return s.repo.Update(client)
}

func (s *clientService) DeleteClient(id uint) error {
	return s.repo.Delete(id)
}

func (s *clientService) GetClient(id uint) (*domain.Client, error) {
	return s.repo.FindByID(id)
}

func (s *clientService) ListClients(page, pageSize int, search string) ([]domain.Client, int64, int64, error) {
	// Sayfalama parametrelerini normalize et
	if page < 1 {
		page = 1
	}
	// Makul bir limit belirle
	if pageSize < 1 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	return s.repo.FindAll(page, pageSize, search)
}

func (s *clientService) GetTotalCount() (int64, error) {
	return s.repo.Count()
}

func (s *clientService) GetClientStatistics() (map[string]int64, error) {
	return s.repo.GetClientStats()
}
