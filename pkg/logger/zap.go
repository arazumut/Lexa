package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// InitLogger, uygulamaya özel yapılandırılmış logger'ı başlatır.
// Environment: "production" ise JSON basar, "development" ise renkli ve okunaklı basar.
func InitLogger(environment string) {
	var config zap.Config

	if environment == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Logları stdout'a (terminale) bas. File yazma opsiyonu sonra eklenebilir.
	logger, err := config.Build()
	if err != nil {
		panic("Logger başlatılamadı: " + err.Error())
	}

	Log = logger
}

// Info, hızlı kullanım için helper
func Info(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Info(msg, fields...)
	}
}

// Error, hızlı kullanım için helper
func Error(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Error(msg, fields...)
	}
}

// Fatal, uygulamayı öldürür
func Fatal(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Fatal(msg, fields...)
		os.Exit(1)
	}
}
