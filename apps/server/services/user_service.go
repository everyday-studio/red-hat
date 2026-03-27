package services

import (
	"fmt"

	"github.com/everyday-studio/redhat/models"
	"github.com/everyday-studio/redhat/repository/postgres"
)

type UserService struct {
	userRepo *postgres.UserRepository
}

func NewUserService(userRepo *postgres.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetByID: %w", err)
	}
	return user, nil
}
