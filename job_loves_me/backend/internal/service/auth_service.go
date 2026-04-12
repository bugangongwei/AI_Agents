package service

import (
	"errors"

	"github.com/bugangongwei/job_loves_me/backend/internal/model"
	"github.com/bugangongwei/job_loves_me/backend/internal/repository"
	"github.com/bugangongwei/job_loves_me/backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(username, email, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := repository.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, *model.User, error) {
	var user model.User
	if err := repository.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}
