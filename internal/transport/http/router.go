package http

import (
	"github.com/arazumut/Lexa/internal/service"
	"github.com/arazumut/Lexa/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

// NewRouter, tüm route tanımlarını ve middleware'leri ayarlar.
func NewRouter(
	r *gin.Engine,
	jwtService service.JWTService,
	authHandler *AuthHandler,
	dashboardHandler *DashboardHandler,
	clientHandler *ClientHandler, // Yeni eklendi
) {
	// 1. PUBLIC ROUTE'LAR (Herkes Girebilir)
	public := r.Group("/")
	{
		public.GET("/login", authHandler.ShowLogin)
		public.POST("/login", authHandler.Login)
		public.GET("/health", HealthCheck)
	}

	// 2. PROTECTED ROUTE'LAR (Sadece Giriş Yapanlar)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtService)) // 🛡️ Kalkan Devrede!
	{
		protected.GET("/", dashboardHandler.Show) // Dashboard
		
		// Müvekkil İşlemleri
		protected.GET("/clients", clientHandler.ShowList)
		protected.GET("/clients/new", clientHandler.ShowCreate)
		protected.GET("/clients/:id/edit", clientHandler.ShowEdit) // Edit Sayfası
		
		protected.GET("/api/clients", clientHandler.List)
		protected.POST("/api/clients", clientHandler.Create)
		protected.PUT("/api/clients/:id", clientHandler.Update)    // Update API
		protected.DELETE("/api/clients/:id", clientHandler.Delete) // Delete API
	}
}
