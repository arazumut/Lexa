package main

import (
	"fmt"
	"log"

	"github.com/arazumut/Lexa/config"
	"github.com/arazumut/Lexa/pkg/database"
)

func main() {
	fmt.Println("⚔️  LEXA: Legal Office Management System Başlatılıyor...")

	// 1. Ayarları Yükle
	cfg := config.LoadConfig()
	fmt.Printf("🔧 Konfigürasyon: Port=%s, Env=%s, DB=%s\n", cfg.AppPort, cfg.Environment, cfg.DBPath)

	// 2. Veritabanına Bağlan
	db, err := database.NewSQLiteDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("❌ Veritabanı hatası: %v", err)
	}
	defer db.Close() // Uygulama kapanırken DB'yi kapat.

	// Şimdilik sadece ayakta kalalım
	log.Println("🚀 Sistem şu an boşta, istek bekleniyor...")
}
