package services

import (
	"errors"
	"strings"

	"campuscore/internal/models"
)

// UserService contains user business logic.
type UserService struct {
	repo models.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(
	repo models.UserRepository,
) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetProfile retrieves a user's profile.
func (s *UserService) GetProfile(
	id string,
) (*models.User, error) {

	id = strings.TrimSpace(id)

	if id == "" {
		return nil, errors.New("user ID is required")
	}

	return s.repo.GetProfile(id)
}

// UpdateProfile updates a user's profile.
func (s *UserService) UpdateProfile(
	user *models.User,
) error {

	if user == nil {
		return errors.New("user is required")
	}

	user.ID = strings.TrimSpace(user.ID)

	if user.ID == "" {
		return errors.New("user ID is required")
	}

	return s.repo.UpdateProfile(user)
}
