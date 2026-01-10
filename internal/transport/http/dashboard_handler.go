package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arazumut/Lexa/internal/domain"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	clientService      domain.ClientService
	caseService        domain.CaseService
	hearingService     domain.HearingService
	transactionService domain.TransactionService
}

func NewDashboardHandler(
	clientService domain.ClientService,
	caseService domain.CaseService,
	hearingService domain.HearingService,
	transactionService domain.TransactionService,
) *DashboardHandler {
	return &DashboardHandler{
		clientService:      clientService,
		caseService:        caseService,
		hearingService:     hearingService,
		transactionService: transactionService,
	}
}

func (h *DashboardHandler) Show(c *gin.Context) {
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	// 1. İstatistikleri Çek
	totalClients, _ := h.clientService.GetTotalCount()
	caseStats, _ := h.caseService.GetCaseStatistics()
	upcomingHearings, _ := h.hearingService.GetUpcomingHearings(5)
	
	// Finansal Veriler (Bakiye ve Grafik)
	balance, monthlyStats, _ := h.transactionService.GetDashboardFinancials()

	// Son Eklenen Dosyalar
	summary, _ := h.caseService.GetDashboardSummary()
	recentCases := summary["recent_cases"]

	// Basit hesaplamalar
	activeCases := caseStats["active"]
	closedCases := caseStats["closed"]
	
	// Para Formatı (Basit)
	totalRevenue := fmt.Sprintf("%.2f₺", balance)

	// Veriyi JSON string'e çevir (JS içinde kullanmak için)
	// Not: Production kodunda hata yönetimi iyi yapılmalı
	statsJSON, _ := json.Marshal(monthlyStats)

	c.HTML(http.StatusOK, "dashboard/dashboard.html", gin.H{
		"title":            "Ana Sayfa - LEXA",
		"email":            email,
		"role":             role,
		"totalClients":     totalClients,
		"activeCases":      activeCases,
		"closedCases":      closedCases,
		"totalRevenue":     totalRevenue,
		"upcomingHearings": upcomingHearings,
		"monthlyStats":     string(statsJSON),
		"recentCases":      recentCases,
	})
	})
}

// GetMiniStats, Layout'un sol barındaki (Sidebar) mini istatistikler için JSON döner.
// Her sayfada (Client, Case vb.) bu endpoint AJAX ile çağrılır.
func (h *DashboardHandler) GetMiniStats(c *gin.Context) {
	// 1. Basitçe countları al (Hata olursa 0 varsayalım)
	totalClients, _ := h.clientService.GetTotalCount()
	caseStats, _ := h.caseService.GetCaseStatistics()
	balance, _, _ := h.transactionService.GetDashboardFinancials()

	// 2. JSON Dön
	c.JSON(http.StatusOK, gin.H{
		"totalFiles":   caseStats["active"] + caseStats["closed"], // Toplam dosya sayısı
		"totalClients": totalClients,
		"totalBalance": fmt.Sprintf("%.2f₺", balance),
	})
}
