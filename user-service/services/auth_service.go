package services

import (
	"errors"
	"user-service/repositories"
	"user-service/utils"
)

type AuthService struct {
	UserRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !utils.ComparePassword(user.PasswordHash, password) {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
