package main

import (
	"fmt"
	"log"

	"github.com/arazumut/Lexa/config"
	"github.com/arazumut/Lexa/internal/repository"
	"github.com/arazumut/Lexa/internal/service"
	"github.com/arazumut/Lexa/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("⚔️  LEXA: Legal Office Management System Başlatılıyor...")

	// 1. Ayarları Yükle
	cfg := config.LoadConfig()
	fmt.Printf("🔧 Konfigürasyon: Port=%s, Env=%s, DB=%s\n", cfg.AppPort, cfg.Environment, cfg.DBPath)

	// 2. Veritabanına Bağlan (GORM)
	db, err := database.NewSQLiteDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("❌ Veritabanı hatası: %v", err)
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
	// internal/transport/http paketini import etmemiz gerekecek, şimdilik inline yapıyorum.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
			"msg":    "Lexa is ready to fight!",
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.String(200, "⚔️ LEXA: Legal Office Management System - AYAKTA!")
	})

	log.Printf("🚀 Sunucu port %s üzerinde başlatılıyor...", cfg.AppPort)
	
	// Uygulamayı başlat ve portu dinle (Bloklayıcı işlem)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("❌ Sunucu başlatılamadı: %v", err)
	}
}
