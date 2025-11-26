package auth

import (
	"errors"
	"go_purple/internal/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepo *user.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepo,
	}
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)

	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}

	_, err = service.UserRepository.Create(user)
	if err != nil {
		return "", err
	}

	return user.Email, nil
}

func (service *AuthService) Login(email, password string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)

	if existedUser == nil {
		return "", errors.New(ErrWrongCredentials)
	}

	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}

	return existedUser.Email, nil
}
