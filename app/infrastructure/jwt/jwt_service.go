package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey = []byte("my_secret_key")

// JWTService интерфейс для работы с JWT токенами
type JWTService interface {
	GenerateToken(email string) (string, error)
	ValidateToken(tokenString string) (*JWTClaim, error)
}

// JWTServiceImpl реализация JWTService
type JWTServiceImpl struct{}

func NewJWTService() JWTService {
	return &JWTServiceImpl{}
}

// JWTClaim структура для хранения данных в JWT токене
type JWTClaim struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken генерирует JWT токен
func (j *JWTServiceImpl) GenerateToken(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &JWTClaim{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateToken проверяет валидность JWT токена
func (j *JWTServiceImpl) ValidateToken(tokenString string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
