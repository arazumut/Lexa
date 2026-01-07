package http

import (
	"github.com/arazumut/Lexa/internal/domain"
	"github.com/arazumut/Lexa/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service domain.UserService
}

// NewAuthHandler, bağımlılıkları enjekte ederek handler'ı oluşturur.
func NewAuthHandler(s domain.UserService) *AuthHandler {
	return &AuthHandler{service: s}
}

// ShowLogin, giriş sayfasını (HTML) render eder.
func (h *AuthHandler) ShowLogin(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{
		"title": "Giriş Yap - LEXA",
	})
}

// Login, formdan gelen verileri işler.
func (h *AuthHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	token, err := h.service.Login(email, password)
	if err != nil {
		logger.Error("Giriş başarısız", zap.String("email", email), zap.Error(err))
		c.HTML(200, "login.html", gin.H{
			"error": "E-posta veya şifre hatalı!",
			"email": email, // Hatalı girişte email silinmesin
		})
		return
	}

	// Başarılı giriş
	// Şimdilik sadece token basıyoruz, ileride Cookie'ye atacağız.
	logger.Info("Kullanıcı giriş yaptı", zap.String("email", email))
	c.String(200, "Giriş Başarılı! Token: "+token)
}
