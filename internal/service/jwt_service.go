package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID uint, email string, role string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey     string
	issuer        string
	expireDuration time.Duration
}

// UserClaims, token içinde taşıyacağımız veriler
type UserClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(secretKey string, issuer string, hours int) JWTService {
	return &jwtService{
		secretKey:      secretKey,
		issuer:         issuer,
		expireDuration: time.Duration(hours) * time.Hour,
	}
}

func (s *jwtService) GenerateToken(userID uint, email string, role string) (string, error) {
	claims := &UserClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expireDuration)),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Algoritma kontrolü (Güvenlik için şart!)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("beklenmeyen imza yöntemi")
		}
		return []byte(s.secretKey), nil
	})
}
