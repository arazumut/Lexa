package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite sürücüsü
)

// NewSQLiteDB, verilen yolda SQLite veritabanı bağlantısı oluşturur.
// Mükemmel bir mimari için connection pooling ayarlarını ve FK desteğini burada yapıyoruz.
func NewSQLiteDB(dbPath string) (*sql.DB, error) {
	// DSN (Data Source Name) ayarları:
	// _foreign_keys=on: İlişkisel bütünlük için ŞART!
	// _busy_timeout=5000: Kilitli DB hatalarını önlemek için 5 sn bekleme.
	dsn := dbPath + "?_foreign_keys=on&_busy_timeout=5000"

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	// Bağlantıyı test et (Ping)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Connection Pooling (Havuz) Ayarları - Performans için kritik!
	// SQLite tek dosya olduğu için çok fazla açık bağlantı "database is locked" hatası verebilir.
	// Ancak WAL modu açarsak limiti artırabiliriz. Şimdilik güvenli limitte tutuyoruz.
	db.SetMaxOpenConns(1) // Write işlemleri için 1 (Eşzamanlı yazma sorunu olmasın)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Veritabanı bağlantısı (SQLite) başarıyla kuruldu:", dbPath)
	return db, nil
}
