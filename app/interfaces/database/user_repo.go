package database

import (
	"database/sql" // Добавляем этот импорт для использования sql.ErrNoRows
	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
	"sharaphka/app/domain"
	"sharaphka/models"
)

type UserRepositoryDB struct {
	db *pop.Connection
}

// NewUserRepositoryDB создает новый репозиторий для работы с пользователями через pop
func NewUserRepositoryDB(db *pop.Connection) *UserRepositoryDB {
	return &UserRepositoryDB{db: db}
}

// FindByEmail ищет пользователя по email в базе данных и возвращает бизнес-объект User
func (r *UserRepositoryDB) FindByEmail(email string) (*domain.User, error) {
	var userModel models.UserModel
	err := r.db.Where("email = ?", email).First(&userModel)
	if err != nil {
		// Проверка на отсутствие записей через sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Пользователь не найден
		}
		return nil, errors.Wrap(err, "failed to find user by email")
	}

	// Преобразуем модель базы данных в бизнес-объект
	userDomain := domain.NewUser(userModel.Email, userModel.Password)

	return userDomain, nil
}
