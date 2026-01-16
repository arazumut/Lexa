package http

import (
	"net/http"

	"github.com/arazumut/Lexa/internal/service"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	service service.SearchService
}

func NewSearchHandler(service service.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

// Search - API Endpoint
// GET /api/search?q=ahmet
func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")

	results, err := h.service.Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Arama hatası"})
		return
	}

	c.JSON(http.StatusOK, results)
}
