package main

import (

	"github.com/arazumut/Lexa/config"
	"github.com/arazumut/Lexa/internal/repository"
	"github.com/arazumut/Lexa/internal/service"
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
	
	// 2. Service (İş Mantığı)
	// Şimdilik kullanılmadığı için alt çizgi (_) ile susturuldu. İleride handler'a verilecek.
	_ = service.NewUserService(userRepo)

	// ---------------------------------------------------------
	// 🌐 HTTP SERVER (WEB KATMANI)
	// ---------------------------------------------------------
	
	// Gin'i release moduna al (Prod ortamı için)
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Basit bir route (Render Health Check için)
	r.GET("/health", func(c *gin.Context) {
		logger.Info("Health check çağrıldı")
		c.JSON(200, gin.H{
			"status": "UP",
			"msg":    "Lexa is ready to fight!",
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.String(200, "⚔️ LEXA: Legal Office Management System - AYAKTA!")
	})

	logger.Info("🚀 Sunucu başlatılıyor...", zap.String("address", ":"+cfg.AppPort))
	
	// Uygulamayı başlat ve portu dinle (Bloklayıcı işlem)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		logger.Fatal("❌ Sunucu başlatılamadı", zap.Error(err))
	}
}
