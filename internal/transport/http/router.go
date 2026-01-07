package http

import (
	"github.com/gin-gonic/gin"
)

// NewRouter, tüm route tanımlarını ve middleware'leri ayarlar.
func NewRouter(r *gin.Engine, authHandler *AuthHandler) {
	// Statik Dosyalar ve Template Ayarları main.go'da yapıldı (veya buraya taşınabilir)
	// Şimdilik route'ları tanımlayalım.

	// Herkese açık route'lar
	public := r.Group("/")
	{
		public.GET("/", authHandler.ShowLogin) // Ana sayfa şimdilik Login olsun
		public.GET("/login", authHandler.ShowLogin)
		public.POST("/login", authHandler.Login)
		public.GET("/health", HealthCheck)
	}

	// Korumalı route'lar (İleride Middleware eklenecek)
	// protected := r.Group("/app")
	// {
	// 	protected.GET("/dashboard", ...)
	// }
}
