package services

import (
	"errors"
	"max/auth/models"
	"max/auth/repositories"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repositories.UserRepository
}

func NewAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) RegisterUser(input models.AuthInput) error {
	// 1. Hashear contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		RoleId:   "2",
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	// 2. Guardar en BD
	return s.repo.CreateUser(user)
}

func (s *AuthService) LoginUser(input models.AuthInput) (string, error) {
	// 1. Buscar usuario
	user, err := s.repo.GetUserByEmail(input.Email)
	if err != nil {
		return "", errors.New("credenciales inválidas")
	}

	// 2. Comparar contraseña (la plana vs la hasheada en BD)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", errors.New("credenciales inválidas")
	}

	// 3. Generar JWT (Payload o "Claims")
	claims := jwt.MapClaims{
		"id":    user.ID,
		"role":  user.RoleId, // Rol del usuario
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Expira en 24 horas
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. Firmar el token con nuestro secreto del .env
	secret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
