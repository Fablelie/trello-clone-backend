package usecase

import (
	"errors"
	"time"

	"github.com/fablelie/trello-clone-backend/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
	secret   string
}

// NewUserUsecase creates a new instance for user business logic
func NewUserUsecase(repo domain.UserRepository, secret string) domain.UserUsecase {
	return &userUsecase{
		userRepo: repo,
		secret:   secret,
	}
}

// Register handles new user registration and password hashing
func (u *userUsecase) Register(email, password, name string) error {
	// Check if the email already exists in the system
	existingUser, _ := u.userRepo.GetByEmail(email)
	if existingUser != nil {
		return errors.New("user already exists with this email")
	}

	// Hash the password before saving to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	// Save the user data to the database via repository
	return u.userRepo.Create(user)
}

// Login verifies password credentials and generates a JWT Token
func (u *userUsecase) Login(email, password string) (string, error) {
	// Find user by email
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Compare the provided password with the hashed one
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Create a JWT Token for authentication
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token หมดอายุใน 72 ชม.
	})

	// Generate encoded token
	return token.SignedString([]byte(u.secret))
}
