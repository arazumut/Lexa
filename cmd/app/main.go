package main

import (

	"github.com/arazumut/Lexa/config"
	"github.com/arazumut/Lexa/internal/repository"
	"github.com/arazumut/Lexa/internal/service"
	transport "github.com/arazumut/Lexa/internal/transport/http" // Alias ile packet adı çakışmasını önle
	"github.com/arazumut/Lexa/pkg/database"
	"github.com/arazumut/Lexa/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. Ayarları Yükle (En Başta!)
	cfg := config.LoadConfig()

	// 2. Logger'ı Başlat (Mükemmel Mimari İçin Şart!)
	logger.InitLogger(cfg.Environment)
	logger.Info("⚔️  LEXA: Legal Office Management System Başlatılıyor...",
		zap.String("env", cfg.Environment),
		zap.String("port", cfg.AppPort),
	)
	
	// Flush: Uygulama kapanırken tüm logları diske/konsola boşaltmayı garanti et.
	defer logger.Log.Sync()

	// 3. Veritabanına Bağlan (GORM)
	db, err := database.NewSQLiteDB(cfg.DBPath)
	if err != nil {
		logger.Fatal("❌ Veritabanı hatası", zap.Error(err))
	}
	
	// GORM'un kendi connection pool yönetimi var ama kapatmak istersek underlying SQL DB'ye erişiriz.
	// main fonksiyonu bitince connection pool da kapanır.
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// ---------------------------------------------------------
	// 🏗️ KATMANLARIN KURULUMU (DEPENDENCY INJECTION)
	// ---------------------------------------------------------
	
	// 1. Repository (Veri Kaynağı)
	userRepo := repository.NewUserRepository(db)
	clientRepo := repository.NewClientRepository(db) // YENİ: Müvekkil Repository
	
	// 2. Service (İş Mantığı)
	// JWT Secret'ı .env'den almalıydık ama şimdilik hardcoded. PROD'da bunu düzeltmeliyiz!
	jwtSecret := "super-secret-key-change-me" 
	jwtService := service.NewJWTService(jwtSecret, "lexa-app", 24) // 24 Saat geçerli
	
	userService := service.NewUserService(userRepo, jwtService)
	clientService := service.NewClientService(clientRepo) // YENİ: Müvekkil Servisi

	// ---------------------------------------------------------
	// ---------------------------------------------------------
	// 🌐 HTTP SERVER (WEB KATMANI)
	// ---------------------------------------------------------
	
	// Gin'i release moduna al (Prod ortamı için)
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Front-end Ayarları
	// custom renderer'ı kullan
	r.HTMLRender = transport.NewRenderer()
	r.Static("/assets", "./web/static/assets")

	// Handler'ları Hazırla
	authHandler := transport.NewAuthHandler(userService)
	dashboardHandler := transport.NewDashboardHandler() 
	clientHandler := transport.NewClientHandler(clientService) // YENİ: Müvekkil Handler
	
	// Router'ı Kur (Dependency Injection)
	// Yeni imzaya uygun olarak clientHandler'ı ekledik.
	transport.NewRouter(r, jwtService, authHandler, dashboardHandler, clientHandler)

	logger.Info("🚀 Sunucu başlatılıyor...", zap.String("address", ":"+cfg.AppPort))
	
	// Uygulamayı başlat ve portu dinle (Bloklayıcı işlem)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		logger.Fatal("❌ Sunucu başlatılamadı", zap.Error(err))
	}
}
