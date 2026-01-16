package http

import (
	"net/http"
	"strconv"

	"github.com/arazumut/Lexa/internal/domain"
	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	service domain.DocumentService
}

func NewDocumentHandler(service domain.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

// Upload - Dosya yükler (API/Form)
func (h *DocumentHandler) Upload(c *gin.Context) {
	// 1. Dosyayı Al
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dosya yüklenemedi: " + err.Error()})
		return
	}

	// 2. Parametreleri Al
	caseIDStr := c.PostForm("case_id")
	category := c.PostForm("category")
	description := c.PostForm("description")

	caseID, err := strconv.Atoi(caseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz Dosya ID"})
		return
	}

	// 3. Kullanıcıyı Al (Context'ten)
	// Middleware userID'yi "userID" key'i ile (uint) set ediyor.
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı bulunamadı"})
		return
	}
	userID := userIDVal.(uint)

	// 4. Servise Gönder
	doc, err := h.service.Upload(file, uint(caseID), userID, category, description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Yükleme başarısız: " + err.Error()})
		return
	}

	// 5. Başarılı Dönüş
	// AJAX ile yükleme yapılıyorsa JSON dön, form submit ise Redirect et.
	// Şimdilik JSON standardı.
	c.JSON(http.StatusOK, gin.H{
		"message": "Dosya başarıyla yüklendi",
		"document": doc,
	})
}

// Delete - Dosya siler
func (h *DocumentHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silme başarısız: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dosya silindi"})
}

// ListByCase - Bir davaya ait dosyaları döner (API)
func (h *DocumentHandler) ListByCase(c *gin.Context) {
	caseID, _ := strconv.Atoi(c.Param("id")) // route: /api/cases/:id/documents

	docs, err := h.service.GetListByCase(uint(caseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosyalar çekilemedi"})
		return
	}

	c.JSON(http.StatusOK, docs)
}
