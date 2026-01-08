package http

import (
	"net/http"

	"github.com/arazumut/Lexa/internal/domain"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	clientService  domain.ClientService
	caseService    domain.CaseService
	hearingService domain.HearingService
}

func NewDashboardHandler(
	clientService domain.ClientService,
	caseService domain.CaseService,
	hearingService domain.HearingService,
) *DashboardHandler {
	return &DashboardHandler{
		clientService:  clientService,
		caseService:    caseService,
		hearingService: hearingService,
	}
}

func (h *DashboardHandler) Show(c *gin.Context) {
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	// 1. İstatistikleri Çek
	totalClients, _ := h.clientService.GetTotalCount()
	caseStats, _ := h.caseService.GetCaseStatistics() // map[string]int64 döner
	upcomingHearings, _ := h.hearingService.GetUpcomingHearings(5)

	// Basit hesaplamalar
	activeCases := caseStats["active"]
	closedCases := caseStats["closed"] // İhtiyaç olursa
	
	// Toplam Tahsilat şu an mock, ileride Accounting modülü ile gerçek olacak
	totalRevenue := "0₺"

	c.HTML(http.StatusOK, "dashboard/dashboard.html", gin.H{
		"title":            "Ana Sayfa - LEXA",
		"email":            email,
		"role":             role,
		"totalClients":     totalClients,
		"activeCases":      activeCases,
		"closedCases":      closedCases,
		"totalRevenue":     totalRevenue,
		"upcomingHearings": upcomingHearings, // Template de range ile dönülecek
	})
}
