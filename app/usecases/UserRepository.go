package usecases

import "sharaphka/app/domain"

// UserRepository — интерфейс для репозитория пользователей
type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
}
