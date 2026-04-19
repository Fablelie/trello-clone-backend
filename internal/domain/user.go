package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	GetByID(id uuid.UUID) (*User, error)
}

type UserUsecase interface {
	Register(email string, password string, name string) error
	Login(email string, password string) (string, error) // return JWT Tokan and error.
}
