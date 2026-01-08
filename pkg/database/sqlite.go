package database

import (
	"log"
	
	"github.com/arazumut/Lexa/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewSQLiteDB, GORM kullanarak SQLite bağlantısı oluşturur.
func NewSQLiteDB(dbPath string) (*gorm.DB, error) {
	// GORM Config
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Development modunda SQL sorgularını gör
	}

	db, err := gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		return nil, err
	}

	// Connection Pooling (GORM üzerinden underlying sql.DB'ye erişerek)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(1) // SQLite için güvenli mod

	log.Println("✅ Veritabanı bağlantısı (GORM + SQLite) başarıyla kuruldu:", dbPath)

	// OTOMATİK MIGRATION (Tabloları struct'lara göre oluşturur)
	// User modelini veritabanına yansıtır.
	// Yeni modeller eklendikçe buraya eklenecek.
	log.Println("📦 Auto-Migration çalıştırılıyor...")
	if err := db.AutoMigrate(&domain.User{}, &domain.Client{}); err != nil {
		return nil, err
	}
	log.Println("✅ Auto-Migration tamamlandı.")

	return db, nil
}
