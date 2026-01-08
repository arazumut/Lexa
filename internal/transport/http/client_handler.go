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

	clients, total, filtered, err := h.service.ListClients(page, pageSize, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veriler çekilemedi"})
		return
	}

	// DataTables Beklenen JSON Formatı
	c.JSON(http.StatusOK, gin.H{
		"draw":            c.Query("draw"), 
		"recordsTotal":    total,    // DB'deki toplam kayıt sayısı
		"recordsFiltered": filtered, // Arama sonrası eşleşen kayıt sayısı
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

// ShowEdit - Müvekkil düzenleme sayfasını render eder.
func (h *ClientHandler) ShowEdit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	client, err := h.service.GetClient(uint(id))
	if err != nil {
		// Müvekkil bulunamazsa listeye yönlendir
		// İdealde 404 sayfası gösterilebilir ama şimdilik redirect daha yumuşak bir geçiş.
		c.Redirect(http.StatusFound, "/clients") 
		return
	}
	
	email, _ := c.Get("email")

	c.HTML(http.StatusOK, "clients/edit.html", gin.H{
		"title": "Müvekkil Düzenle - LEXA",
		"email": email,
		"client": client, // Mevcut veriyi template'e gönder
	})
}

// Update - Müvekkil bilgilerini günceller (API Endpoint).
func (h *ClientHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var client domain.Client
	
	// JSON verisini struct'a bind et
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri formatı"})
		return
	}
	
	// URL'den gelen ID'yi güvenli bir şekilde struct'a ata.
	// Kullanıcının JSON body içinde manipüle etmesini engelleriz.
	client.ID = uint(id)
	
	if err := h.service.UpdateClient(&client); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Müvekkil başarıyla güncellendi",
	})
}

// Delete - Müvekkili siler (API Endpoint).
func (h *ClientHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.service.DeleteClient(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silme işlemi başarısız: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Müvekkil başarıyla silindi",
	})
}
