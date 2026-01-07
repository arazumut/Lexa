# Build Stage
FROM golang:1.25rc1-alpine3.21 AS builder

# SQLite için CGO gerekli (gcc, musl-dev)
RUN apk add --no-cache build-base

WORKDIR /app

# Bağımlılıkları kopyala ve indir
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodları kopyala
COPY . .

# Uygulamayı derle (CGO_ENABLED=1 SQLite için şart)
# -ldflags "-s -w" binary boyutunu küçültür
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-s -w" -o lexa cmd/app/main.go

# Runtime Stage
FROM alpine:latest

WORKDIR /app

# Gerekli sertifikaları ekle (HTTPS istekleri için)
RUN apk add --no-cache ca-certificates

# Builder aşamasından binary'i al
COPY --from=builder /app/lexa .

# Web şablonlarını ve statik dosyaları kopyala
COPY --from=builder /app/web ./web

# Portu dışarı aç
EXPOSE 8080

# Uygulamayı başlat
CMD ["./lexa"]
