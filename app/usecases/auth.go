package usecases

import (
	"errors"
	"sharaphka/app/infrastructure/jwt"
)

// AuthUseCase отвечает за логику аутентификации
type AuthUseCase struct {
	userRepo   UserRepository
	jwtService jwt.JWTService
}

// NewAuthUseCase создаёт новый AuthUseCase
func NewAuthUseCase(userRepo UserRepository, jwtService jwt.JWTService) *AuthUseCase {
	return &AuthUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Login выполняет аутентификацию пользователя и возвращает JWT токен
func (a *AuthUseCase) Login(email, password string) (string, error) {
	user, err := a.userRepo.FindByEmail(email)
	if err != nil || user == nil || user.Password != password {
		return "", errors.New("invalid credentials")
	}

	// Генерация JWT токена
	token, err := a.jwtService.GenerateToken(email)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateToken проверяет валидность JWT токена
func (a *AuthUseCase) ValidateToken(token string) (bool, error) {
	claims, err := a.jwtService.ValidateToken(token)
	if err != nil {
		return false, err
	}

	// Здесь можно добавить дополнительную логику, например, проверку срока действия токена
	return claims != nil, nil
}
