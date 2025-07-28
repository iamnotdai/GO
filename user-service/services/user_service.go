package services

import (
	"user-service/dtos"
	"user-service/models"
	"user-service/repositories"
	"user-service/utils"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) CreateUser(input dtos.CreateUserInput) error {
	// Hash mật khẩu
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		Role:         input.Role,
		PasswordHash: hashedPassword,
	}

	return s.Repo.CreateUser(&user)
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAllUsers()
}
