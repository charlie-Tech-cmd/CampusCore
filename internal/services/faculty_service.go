package services

import (
	"campuscore/internal/models"
)

// FacultyService contains faculty business logic.
type FacultyService struct {
	repo models.FacultyRepository
}

// NewFacultyService creates a FacultyService.
func NewFacultyService(
	repo models.FacultyRepository,
) *FacultyService {
	return &FacultyService{
		repo: repo,
	}
}

// CreateFaculty creates a faculty.
func (s *FacultyService) CreateFaculty(faculty *models.Faculty) error {
	return s.repo.Create(faculty)
}

// GetFaculty retrieves a faculty by ID.
func (s *FacultyService) GetFaculty(id int) (*models.Faculty, error) {
	return s.repo.FindByID(id)
}

// GetFacultyByCode retrieves a faculty by code.
func (s *FacultyService) GetFacultyByCode(code string) (*models.Faculty, error) {
	return s.repo.FindByCode(code)
}

// ListFaculties retrieves all faculties.
func (s *FacultyService) ListFaculties() ([]models.Faculty, error) {
	return s.repo.List()
}

// UpdateFaculty updates a faculty.
func (s *FacultyService) UpdateFaculty(faculty *models.Faculty) error {
	return s.repo.Update(faculty)
}

// DeleteFaculty removes a faculty.
func (s *FacultyService) DeleteFaculty(id int) error {
	return s.repo.Delete(id)
}
