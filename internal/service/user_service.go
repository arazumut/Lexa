package service

import (
	"errors"
	"time"

	"github.com/arazumut/Lexa/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo       domain.UserRepository
	jwtService JWTService
}

// NewUserService, iş mantığını içeren servisi oluşturur.
func NewUserService(repo domain.UserRepository, jwtService JWTService) domain.UserService {
	return &userService{
		repo:       repo,
		jwtService: jwtService,
	}
}

// Register, yeni kullanıcı kaydeder. Şifreyi hash'ler.
func (s *userService) Register(email, password, name string) error {
	// 1. Kullanıcı zaten var mı?
	existingUser, err := s.repo.GetByEmail(email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("bu email adresi zaten kayıtlı")
	}

	// 2. Şifreyi Hashle (Cost 10 ideal)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 3. Kullanıcıyı Oluştur
	user := &domain.User{
		Email:     email,
		Password:  string(hashedPassword),
		Name:      name,
		Role:      "staff", // Default
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.Create(user)
}

// Login, email ve şifre kontrolü yapar.
// Başarılı ise JWT token döner.
func (s *userService) Login(email, password string) (string, error) {
	// 1. Kullanıcıyı bul
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("kullanıcı bulunamadı veya şifre hatalı")
	}

	// 2. Şifreyi doğrula
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("kullanıcı bulunamadı veya şifre hatalı") // Güvenlik için genel hata mesajı
	}

	// 3. Token oluştur (Artık GERÇEK JWT!)
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) CreateUser(user *domain.User) error {
	// 1. Kullanıcı zaten var mı?
	existingUser, err := s.repo.GetByEmail(user.Email)
	if err != nil {
		// Repo hatası (DB bağlantısı vb.)
		// Not: GetByEmail nil, nil dönüyorsa sorun yok.
		// Ama nil, err dönüyorsa hata var.
		// Burada GetByEmail implementasyonuna güveniyoruz.
		return err
	}
	if existingUser != nil {
		return errors.New("bu email adresi zaten kayıtlı")
	}

	// 2. Şifreyi Hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Zaman damgaları
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	user.UpdatedAt = time.Now()

	// 3. Veritabanı Kaydı
	return s.repo.Create(user)
}
