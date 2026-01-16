package service

import (
	"fmt"
	"sync"

	"github.com/arazumut/Lexa/internal/domain"
)

// SearchResult, arama sonucunda dönecek ortak yapı
type SearchResult struct {
	Type  string `json:"type"`  // "client", "case"
	ID    uint   `json:"id"`
	Title string `json:"title"` // Görünen isim (Ad Soyad veya Dava Başlığı)
	Meta  string `json:"meta"`  // Alt bilgi (TCKN veya Dosya No)
	URL   string `json:"url"`   // Tıklanınca gideceği yer
}

type SearchService interface {
	Search(query string) ([]SearchResult, error)
}

type searchService struct {
	clientRepo domain.ClientRepository
	caseRepo   domain.CaseRepository
}

func NewSearchService(clientRepo domain.ClientRepository, caseRepo domain.CaseRepository) SearchService {
	return &searchService{
		clientRepo: clientRepo,
		caseRepo:   caseRepo,
	}
}

// Search - Hem Client hem Case tablolarında arama yapar ve sonuçları birleştirir.
// Performans için Goroutine (Concurrency) kullanılabilir.
func (s *searchService) Search(query string) ([]SearchResult, error) {
	var results []SearchResult
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Aramak için en az 2 karakter girilsin
	if len(query) < 2 {
		return []SearchResult{}, nil
	}

	wg.Add(2)

	// 1. Müvekkil Araması (Goroutine)
	go func() {
		defer wg.Done()
		// Repository'de search fonksiyonu zaten vardı (List içinde), ama burada ham dataya ihtiyacımız var.
		// Hız için repo metodunu kullanıp, dönüştüreceğiz.
		// Not: ListClients biraz ağır olabilir (Count hesaplıyor), 
		// sadece arama için repository'e "Search" metodu eklemek daha doğru olur ama
		// şimdilik mevcut "FindAll" metodunu limitli kullanalım.
		clients, _, _, err := s.clientRepo.FindAll(1, 5, query) 
		if err == nil {
			mu.Lock()
			for _, c := range clients {
				results = append(results, SearchResult{
					Type:  "Müvekkil",
					ID:    c.ID,
					Title: c.Name,
					Meta:  c.Identity,
					URL:   fmt.Sprintf("/clients/%d", c.ID),
				})
			}
			mu.Unlock()
		}
	}()

	// 2. Dava Araması (Goroutine)
	go func() {
		defer wg.Done()
		// CaseRepo FindAll metodunu kullanalım
		cases, _, _, err := s.caseRepo.FindAll(1, 5, query, 0)
		if err == nil {
			mu.Lock()
			for _, c := range cases {
				results = append(results, SearchResult{
					Type:  "Dosya",
					ID:    c.ID,
					Title: c.FileNumber + " - " + c.Title, // Örn: 2024/1 E. - Boşanma Davası
					Meta:  c.Court,
					URL:   fmt.Sprintf("/cases/%d", c.ID),
				})
			}
			mu.Unlock()
		}
	}()

	wg.Wait()

	return results, nil
}
