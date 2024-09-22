package domain

// User представляет бизнес-логику пользователя
type User struct {
	Email    string
	Password string
}

// NewUser создает новый бизнес-объект пользователя
func NewUser(email, password string) *User {
	return &User{
		Email:    email,
		Password: password,
	}
}
