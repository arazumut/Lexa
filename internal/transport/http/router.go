package http

import (
	"github.com/arazumut/Lexa/internal/service"
	"github.com/arazumut/Lexa/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

// NewRouter, tüm route tanımlarını ve middleware'leri ayarlar.
func NewRouter(r *gin.Engine, jwtService service.JWTService, authHandler *AuthHandler, dashboardHandler *DashboardHandler) {
	// Statik Route'lar (main.go'da tanımlıydı ama burası daha temiz olurdu, neyse)

	// 1. PUBLIC ROUTE'LAR (Herkes Girebilir)
	public := r.Group("/")
	{
		public.GET("/login", authHandler.ShowLogin)
		public.POST("/login", authHandler.Login)
		public.GET("/health", HealthCheck)
		// Register sayfası ileride eklenebilir
	}

	// 2. PROTECTED ROUTE'LAR (Sadece Giriş Yapanlar)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtService)) // 🛡️ Kalkan Devrede!
	{
		protected.GET("/", dashboardHandler.Show) // Ana Sayfa artık Dashboard
		// İleride /clients, /cases gibi yollar buraya gelecek
	}
}
