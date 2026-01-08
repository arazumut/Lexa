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
	clientHandler *ClientHandler,
	caseHandler *CaseHandler,
	hearingHandler *HearingHandler, // Yeni eklendi
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
		protected.GET("/clients/:id/edit", clientHandler.ShowEdit)
		
		protected.GET("/api/clients", clientHandler.List)
		protected.POST("/api/clients", clientHandler.Create)
		protected.PUT("/api/clients/:id", clientHandler.Update)
		protected.DELETE("/api/clients/:id", clientHandler.Delete)

		// Dava (Case) İşlemleri
		protected.GET("/cases", caseHandler.ShowList)
		protected.GET("/cases/new", caseHandler.ShowCreate)
		protected.GET("/cases/:id/edit", caseHandler.ShowEdit)
		
		protected.GET("/api/cases", caseHandler.List)
		protected.POST("/api/cases", caseHandler.Create)
		protected.PUT("/api/cases/:id", caseHandler.Update)
		protected.DELETE("/api/cases/:id", caseHandler.Delete)

		// Duruşma (Hearing) İşlemleri
		protected.GET("/hearings", hearingHandler.ShowList)
		protected.GET("/hearings/new", hearingHandler.ShowCreate)
		protected.GET("/hearings/:id/edit", hearingHandler.ShowEdit)
		
		protected.GET("/api/hearings", hearingHandler.List)
		protected.POST("/api/hearings", hearingHandler.Create)
		protected.PUT("/api/hearings/:id", hearingHandler.Update)
		protected.DELETE("/api/hearings/:id", hearingHandler.Delete)
	}
}
