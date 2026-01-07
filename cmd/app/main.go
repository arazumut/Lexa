package main

import (
	"fmt"
	"log"

	"github.com/arazumut/Lexa/config"
	"github.com/arazumut/Lexa/internal/repository"
	"github.com/arazumut/Lexa/internal/service"
	"github.com/arazumut/Lexa/pkg/database"
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
	userService := service.NewUserService(userRepo)

	// ---------------------------------------------------------
	// 🧪 HIZLI TEST (DEBUG İÇİN - SİLİNECEK)
	// ---------------------------------------------------------
	log.Println("🧪 'admin@lexa.com' kullanıcısı oluşturuluyor (Test)...")
	err = userService.Register("admin@lexa.com", "123456", "Sistem Yöneticisi")
	if err != nil {
		log.Printf("⚠️ Kullanıcı oluşturma uyarısı: %v", err)
	} else {
		log.Println("✅ Test kullanıcısı başarıyla oluşturuldu!")
	}
	
	// Login Testi
	token, err := userService.Login("admin@lexa.com", "123456")
	if err != nil {
		log.Printf("❌ Login başarısız: %v", err)
	} else {
		log.Printf("✅ Login başarılı! Token: %s", token)
	}

	// Şimdilik sadece ayakta kalalım
	log.Println("🚀 Sistem şu an boşta, istek bekleniyor...")
}
