package http

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/arazumut/Lexa/internal/domain"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service       domain.TransactionService
	clientService domain.ClientService
	caseService   domain.CaseService
}

func NewTransactionHandler(
	service domain.TransactionService,
	clientService domain.ClientService,
	caseService domain.CaseService,
) *TransactionHandler {
	return &TransactionHandler{
		service:       service,
		clientService: clientService,
		caseService:   caseService,
	}
}

// ShowList - Muhasebe listesi
func (h *TransactionHandler) ShowList(c *gin.Context) {
	email, _ := c.Get("email")
	
	// Toplam bakiyeyi de sayfada göstermek güzel olur
	balance, _, _ := h.service.GetDashboardFinancials()
	formattedBalance := fmt.Sprintf("%.2f₺", balance)

	c.HTML(http.StatusOK, "accounting/list.html", gin.H{
		"title":   "Muhasebe & Finans - LEXA",
		"email":   email,
		"balance": formattedBalance,
	})
}

// ShowCreate - Yeni İşlem Ekleme Sayfası
func (h *TransactionHandler) ShowCreate(c *gin.Context) {
	email, _ := c.Get("email")

	// Dropdownlar için verileri çek
	// Not: Performans için ileride Select2 AJAX kullanılabilir ama şimdilik tümünü çekiyoruz.
	clients, _, _, _ := h.clientService.ListClients(1, 1000, "")
	cases, _, _, _ := h.caseService.ListCases(1, 1000, "", 0)

	// URL'den gelen parametreler (Örn: Dava detayından "Tahsilat Ekle" denirse)
	preSelectedCaseID := c.Query("case_id")
	preSelectedClientID := c.Query("client_id")
	preSelectedType := c.Query("type") // income veya expense

	c.HTML(http.StatusOK, "accounting/create.html", gin.H{
		"title":      "Yeni Finansal İşlem - LEXA",
		"email":      email,
		"clients":    clients,
		"cases":      cases,
		"selCaseID":  preSelectedCaseID,
		"selClientID": preSelectedClientID,
		"selType":    preSelectedType,
	})
}

// List - API Endpoint (DataTables)
func (h *TransactionHandler) List(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("length", "10"))
	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))
	// search := c.Query("search[value]") // Filtreleme eklenecek
	
	page := (start / pageSize) + 1

	filter := domain.TransactionFilter{
		// Search: search, // Repo'da implemente etmiştik
	}

	transactions, total, err := h.service.ListTransactions(page, pageSize, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veriler çekilemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"draw":            c.Query("draw"),
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            transactions,
	})
}

// Create - API Endpoint
func (h *TransactionHandler) Create(c *gin.Context) {
	var t domain.Transaction
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri: " + err.Error()})
		return
	}

	// Tarih kontrolü (Repo/Service yapıyor ama burada da parse hatası olmasın)
	// JSON binding time.Time'ı RFC3339 bekler.

	if err := h.service.CreateTransaction(&t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "İşlem kaydedildi", "id": t.ID})
}

// Delete - API Endpoint
func (h *TransactionHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.DeleteTransaction(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silinemedi: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "İşlem silindi"})
}
