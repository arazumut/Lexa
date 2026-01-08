package http

import (
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) Show(c *gin.Context) {
	// Middleware sayesinde userID ve role context'te mevcut
	email, _ := c.Get("email")     // Claims'e email eklemiştik? Bakalım jwt_service.go'ya.
	role, _ := c.Get("role")

	c.HTML(200, "dashboard.html", gin.H{
		"title": "Ana Sayfa - LEXA",
		"email": email,
		"role":  role,
	})
}
