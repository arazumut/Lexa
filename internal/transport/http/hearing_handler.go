package http

import (
	"net/http"
	"strconv"

	"github.com/arazumut/Lexa/internal/domain"
	"github.com/gin-gonic/gin"
)

type HearingHandler struct {
	service     domain.HearingService
	caseService domain.CaseService // Dropdown doldurmak için davelara ihtiyacımız var
}

func NewHearingHandler(service domain.HearingService, caseService domain.CaseService) *HearingHandler {
	return &HearingHandler{
		service:     service,
		caseService: caseService,
	}
}

// ShowList - Duruşma takvimi/listesi sayfasını render eder.
func (h *HearingHandler) ShowList(c *gin.Context) {
	email, _ := c.Get("email")
	
	c.HTML(http.StatusOK, "hearings/list.html", gin.H{
		"title": "Duruşma Takvimi - LEXA",
		"email": email,
	})
}

// ShowCreate - Yeni duruşma ekleme sayfasını render eder.
func (h *HearingHandler) ShowCreate(c *gin.Context) {
	email, _ := c.Get("email")

	// Aktif davaları getir (Selectbox için)
	// TODO: İleride sadece 'active' statüsündeki davaları çeken özel bir metod yazılabilir.
	// Şimdilik ListCases kullanıyoruz.
	cases, _, _, _ := h.caseService.ListCases(1, 200, "", 0)

	// URL'den pre-select için case_id gelebilir (Dava detayından "Duruşma Ekle" denirse)
	selectedCaseID := c.Query("case_id")

	c.HTML(http.StatusOK, "hearings/create.html", gin.H{
		"title":          "Yeni Duruşma Ekle - LEXA",
		"email":          email,
		"cases":          cases,
		"selectedCaseID": selectedCaseID,
	})
}

// ShowEdit - Duzenleme Sayfasi
func (h *HearingHandler) ShowEdit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	hearing, err := h.service.GetHearing(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/hearings")
		return
	}
	
	email, _ := c.Get("email")
	// Dropdown için yine davalari çek
	cases, _, _, _ := h.caseService.ListCases(1, 200, "", 0)

	c.HTML(http.StatusOK, "hearings/edit.html", gin.H{
		"title":   "Duruşma Düzenle - LEXA",
		"email":   email,
		"hearing": hearing,
		"cases":   cases,
	})
}

// List - API Endpoint (DataTables veya Calendar için JSON döner)
func (h *HearingHandler) List(c *gin.Context) {
	// DataTables Standart Parametreleri (İleride Calendar.js formatına dönüştürülebilir)
	// Basitlik için şimdilik DataTables mantığıyla ilerliyoruz.
	// Takvim görünümü için "FullCalendar.io" entegrasyonu yapılabilir, o zaman formatı değiştiririz.
	
	pageSize, _ := strconv.Atoi(c.DefaultQuery("length", "10"))
	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))
	// search := c.Query("search[value]") // Repository'de henüz search yok, ekleriz.
	
	page := (start / pageSize) + 1

	hearings, total, err := h.service.ListHearings(page, pageSize, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veriler çekilemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"draw":            c.Query("draw"), 
		"recordsTotal":    total,
		"recordsFiltered": total, // Search implemente edilince düzeltilmeli
		"data":            hearings,
	})
}

// CalendarEvents - FullCalendar için JSON formatında veri döner
func (h *HearingHandler) CalendarEvents(c *gin.Context) {
	// Takvim görünümü için daha geniş bir aralık çekebiliriz.
	// Şimdilik son 1000 duruşmayı çekelim. 
	// İdeal olan start/end parametrelerine göre filtrelemektir.
	hearings, _, err := h.service.ListHearings(1, 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veriler çekilemedi"})
		return
	}

	var events []gin.H
	for _, hearing := range hearings {
		color := "bg-info"
		if hearing.Status == "completed" {
			color = "bg-success"
		} else if hearing.Status == "cancelled" {
			color = "bg-danger"
		} else if hearing.Status == "postponed" {
			color = "bg-warning"
		}

		title := hearing.Title
		if hearing.Case.ID != 0 {
			title = hearing.Case.FileNumber + " - " + title
		}

		events = append(events, gin.H{
			"id":        hearing.ID,
			"title":     title,
			"start":     hearing.Date, // ISO8601
			"className": color,
			"url":       "/hearings/" + strconv.Itoa(int(hearing.ID)) + "/edit", // Düzenlemeye git
		})
	}

	c.JSON(http.StatusOK, events)
}

func (h *HearingHandler) Create(c *gin.Context) {
	var hearing domain.Hearing
	
	// Tarih formatı RFC3339 olarak gelmeli (JS: toISOString())
	if err := c.ShouldBindJSON(&hearing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri: " + err.Error()})
		return
	}
	
	// Tarih kontrolü (Eski tarihe duruşma eklenmemeli uyarısı? Yoksa izin verilebilir, geçmişi kaydetmek için)
	// Şimdilik izin veriyoruz.

	if err := h.service.CreateHearing(&hearing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kaydedilemedi: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Duruşma kaydedildi",
		"id":      hearing.ID,
	})
}

// Update - API Endpoint
func (h *HearingHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var hearing domain.Hearing
	if err := c.ShouldBindJSON(&hearing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}
	hearing.ID = uint(id)
	
	if err := h.service.UpdateHearing(&hearing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Duruşma güncellendi",
	})
}

// Delete - API Endpoint
func (h *HearingHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.DeleteHearing(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silinemedi: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Duruşma silindi",
	})
}
