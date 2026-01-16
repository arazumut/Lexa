package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/arazumut/Lexa/internal/domain"
	"github.com/gin-gonic/gin"
)

type CaseHandler struct {
	service       domain.CaseService
	clientService domain.ClientService
}

// NewCaseHandler constructor
// UI'da dropdown doldurmak için ClientService'e de ihtiyacımız var.
func NewCaseHandler(service domain.CaseService, clientService domain.ClientService) *CaseHandler {
	return &CaseHandler{
		service:       service,
		clientService: clientService,
	}
}

// ShowList - Dava listesi sayfasını render eder.
func (h *CaseHandler) ShowList(c *gin.Context) {
	email, _ := c.Get("email")
	
	c.HTML(http.StatusOK, "cases/list.html", gin.H{
		"title": "Dava Dosyaları - LEXA",
		"email": email,
	})
}

// ShowCreate - Yeni dava ekleme sayfasını render eder.
func (h *CaseHandler) ShowCreate(c *gin.Context) {
	email, _ := c.Get("email")

	// Müvekkil seçimi için tüm müvekkilleri getir (basit bir liste, sayfalama olmadan)
	// Not: Prodüksiyonda binlerce müvekkil varsa bu dropdown AJAX Select2 ile yapılmalı.
	// Şimdilik 100 müvekkil limitiyle listeyi alalım.
	clients, _, _, _ := h.clientService.ListClients(1, 100, "")

	// URL'den gelen pre-select parametresi
	selectedClientID := c.Query("client_id")

	c.HTML(http.StatusOK, "cases/create.html", gin.H{
		"title":            "Yeni Dava Aç - LEXA",
		"email":            email,
		"clients":          clients, // Dropdown için
		"selectedClientID": selectedClientID,
	})
}

// ShowEdit - Dava düzenleme sayfasını render eder.
func (h *CaseHandler) ShowEdit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	caseFile, err := h.service.GetCase(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cases")
		return
	}
	
	email, _ := c.Get("email")
	// Mevcut müvekkilini de listede göstermek için client listesi lazım
	clients, _, _, _ := h.clientService.ListClients(1, 100, "")

	c.HTML(http.StatusOK, "cases/edit.html", gin.H{
		"title":   "Dava Düzenle - LEXA",
		"email":   email,
		"case":    caseFile,
		"clients": clients,
	})
}

// ShowDetail - Dava dosyasının tüm detaylarını gösterir (Evraklar, Duruşmalar, Muhasebe dahil)
func (h *CaseHandler) ShowDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// Repository'den Preload ile her şeyi çekiyoruz.
	caseFile, err := h.service.GetCase(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cases")
		return
	}

	email, _ := c.Get("email")
	
	c.HTML(http.StatusOK, "cases/detail.html", gin.H{
		"title": caseFile.FileNumber + " - Dava Detayı",
		"email": email,
		"case":  caseFile, 
		// "case.Documents", "case.Hearings", "case.Transactions" zaten içinde dolu geliyor.
	})
}


// List - DataTables API Endpoint
func (h *CaseHandler) List(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("length", "10"))
	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))
	search := c.Query("search[value]")
	
	page := (start / pageSize) + 1

	cases, total, filtered, err := h.service.ListCases(page, pageSize, search, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veriler çekilemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"draw":            c.Query("draw"), 
		"recordsTotal":    total,
		"recordsFiltered": filtered,
		"data":            cases,
	})
}

// Create - API Endpoint
func (h *CaseHandler) Create(c *gin.Context) {
	// Gelen JSON'u karşılayacak geçici bir struct (Tarih parsing için gerekebilir)
	// Ancak time.Time standard JSON unmarshal destekler (RFC3339).
	var caseReq domain.Case
	
	if err := c.ShouldBindJSON(&caseReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri: " + err.Error()})
		return
	}

	// Otomatik set edilmesi gereken alanlar varsa frontend göndermemeli, burada set edilmeli.
	// Örneğin StartDate boşsa bugünün tarihi atanabilir:
	if caseReq.StartDate.IsZero() {
		caseReq.StartDate = time.Now()
	}

	if err := h.service.CreateCase(&caseReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kaydedilemedi: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dava dosyası başarıyla oluşturuldu",
		"id":      caseReq.ID,
	})
}

// Update - API Endpoint
func (h *CaseHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var caseReq domain.Case
	if err := c.ShouldBindJSON(&caseReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}
	
	caseReq.ID = uint(id)
	
	if err := h.service.UpdateCase(&caseReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Dava dosyası güncellendi",
	})
}

// Delete - API Endpoint
func (h *CaseHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.service.DeleteCase(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silinemedi: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Dava dosyası silindi",
	})
}
