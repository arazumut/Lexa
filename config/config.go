package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      string
	DBPath       string
	Environment  string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("UYARI: .env dosyası bulunamadı, ortam değişkenleri kontrol ediliyor...")
	}

	return &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DBPath:      getEnv("DB_PATH", "./lexa.db"),
		Environment: getEnv("ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
