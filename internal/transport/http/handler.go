package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Sağlık kontrolü (Health Check) handler'ı.
// Render deploy'unun başarılı olduğunu anlamak için bunu kullanacağız.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "alive",
		"message": "⚔️ LEXA System is Running!",
	})
}
