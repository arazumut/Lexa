.PHONY: run build clean

APP_NAME=lexa
MAIN_FILE=cmd/app/main.go

run:
	@echo "🚀 Uygulama başlatılıyor..."
	@go run $(MAIN_FILE)

build:
	@echo "🔨 Derleniyor..."
	@go build -o bin/$(APP_NAME) $(MAIN_FILE)

clean:
	@echo "🧹 Temizlik yapılıyor..."
	@rm -rf bin/
