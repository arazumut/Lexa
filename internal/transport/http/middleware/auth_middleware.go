package middleware

import (
	"net/http"
	"strings"

	"github.com/arazumut/Lexa/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware, korumalı route'lar için kimlik doğrulaması yapar.
func AuthMiddleware(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Token'ı Cookie'den almaya çalış
		tokenString, err := c.Cookie("Authorization")
		
		// 2. Cookie yoksa Header'a bak (Bearer Token)
		if err != nil {
			// Debug Log
			// logger.Info("Cookie bulunamadı, Header kontrol ediliyor...", zap.Error(err))
			
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				// Token hiç yok
				// logger.Warn("🚫 Yetkisiz Giriş: Token yok, Login'e yönlendiriliyor.")
				c.Redirect(http.StatusFound, "/login")
				c.Abort()
				return
			}
			// "Bearer " kısmını temizle
			tokenString = strings.Replace(authHeader, "Bearer ", "", 1)
		}

		// 3. Token'ı Doğrula
		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			// logger.Warn("🚫 Geçersiz Token", zap.Error(err))
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// 4. Token içindeki bilgileri Context'e yükle (ki diğer handler'lar kullanabilsin)
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			// UserID float64 olarak döner JSON parse işlemlerinde, uint'e çeviriyoruz
			if uid, ok := claims["user_id"].(float64); ok {
				c.Set("userID", uint(uid))
			}
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}
		}

		c.Next()
	}
}
