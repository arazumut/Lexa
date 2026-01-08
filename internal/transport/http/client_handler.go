package http

import (
	"net/http"
	"strconv"

	"github.com/arazumut/Lexa/internal/domain"
	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	service domain.ClientService
}

func NewClientHandler(service domain.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

// ShowList - Müvekkil listesi sayfasını render eder.
func (h *ClientHandler) ShowList(c *gin.Context) {
	// Kullanıcı bilgilerini al
	email, _ := c.Get("email")
	
	c.HTML(http.StatusOK, "clients/list.html", gin.H{
		"title": "Müvekkiller - LEXA",
		"email": email, // Base template için gerekli
	})
}

// ShowCreate - Yeni müvekkil ekleme sayfasını render eder.
func (h *ClientHandler) ShowCreate(c *gin.Context) {
	email, _ := c.Get("email")

	c.HTML(http.StatusOK, "clients/create.html", gin.H{
		"title": "Yeni Müvekkil Ekle - LEXA",
		"email": email,
	})
}

// List - DataTables için JSON veri döner.
func (h *ClientHandler) List(c *gin.Context) {
	// DataTables Parametreleri
	// draw: Güvenlik/Sıra sayacı
	// start: Başlangıç kaydı (offset)
	// length: Sayfa başına kayıt sayısı (limit)
	// search[value]: Arama terimi
	
	pageSize, _ := strconv.Atoi(c.DefaultQuery("length", "10"))
	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))
	search := c.Query("search[value]")
	
	// Start -> Page Dönüşümü
	page := (start / pageSize) + 1

	clients, total, err := h.service.ListClients(page, pageSize, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veriler çekilemedi"})
		return
	}

	// DataTables Beklenen JSON Formatı
	c.JSON(http.StatusOK, gin.H{
		"draw":            c.Query("draw"), 
		"recordsTotal":    total, // Filtresiz toplam sayıyı da buraya verebiliriz aslında ama şimdilik aynı
		"recordsFiltered": total, // Filtrelenmiş toplam sayı
		"data":            clients,
	})
}

// Create - Yeni müvekkil kaydeder (API Endpoint).
func (h *ClientHandler) Create(c *gin.Context) {
	var client domain.Client
	
	// JSON verisini struct'a bind et
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri formatı: " + err.Error()})
		return
	}

	// Servis katmanına ilet
	if err := h.service.CreateClient(&client); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kaydedilemedi: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Müvekkil başarıyla oluşturuldu",
		"id":      client.ID,
	})
}
