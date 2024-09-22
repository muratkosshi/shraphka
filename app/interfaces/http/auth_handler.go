package http

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"net/http"
	"sharaphka/app/usecases"
)

type AuthHandler struct {
	authUseCase *usecases.AuthUseCase
}

func NewAuthHandler(authUseCase *usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

// Login выполняет аутентификацию пользователя
func (h *AuthHandler) Login(c buffalo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{"error": "Invalid request"}))
	}

	token, err := h.authUseCase.Login(req.Email, req.Password)
	if err != nil {
		return c.Render(http.StatusUnauthorized, render.JSON(map[string]string{"error": "Invalid credentials"}))
	}

	return c.Render(http.StatusOK, render.JSON(map[string]string{"token": token}))
}

// AuthMiddleware проверяет JWT токен
func (h *AuthHandler) AuthMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// Логика проверки токена
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.Render(http.StatusUnauthorized, render.JSON(map[string]string{"error": "Missing token"}))
		}

		// Валидация токена через use case
		_, err := h.authUseCase.ValidateToken(token)
		if err != nil {
			return c.Render(http.StatusUnauthorized, render.JSON(map[string]string{"error": "Invalid token"}))
		}

		// Продолжить выполнение следующего обработчика
		return next(c)
	}
}

// Profile возвращает профиль пользователя
func (h *AuthHandler) Profile(c buffalo.Context) error {
	// Пример профиля
	return c.Render(http.StatusOK, render.JSON(map[string]string{"profile": "User profile data"}))
}
