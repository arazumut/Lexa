package http

import (
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

	// Basit hesaplamalar
	activeCases := caseStats["active"]
	closedCases := caseStats["closed"]
	
	// Para Formatı (Basit)
	totalRevenue := fmt.Sprintf("%.2f₺", balance)

	c.HTML(http.StatusOK, "dashboard/dashboard.html", gin.H{
		"title":            "Ana Sayfa - LEXA",
		"email":            email,
		"role":             role,
		"totalClients":     totalClients,
		"activeCases":      activeCases,
		"closedCases":      closedCases,
		"totalRevenue":     totalRevenue,
		"upcomingHearings": upcomingHearings,
		"monthlyStats":     monthlyStats, // Grafik için (JS tarafında json olarak kullanılacak)
	})
}
