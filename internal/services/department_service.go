package services

import (
	"campuscore/internal/models"
)

// DepartmentRepository defines the persistence contract.
type DepartmentRepository interface {
	Create(*models.Department) error
	FindByID(int) (*models.Department, error)
	FindByCode(string) (*models.Department, error)
	List() ([]models.Department, error)
	Update(*models.Department) error
	Delete(int) error
}

// DepartmentService contains department business logic.
type DepartmentService struct {
	repo DepartmentRepository
}

// NewDepartmentService creates a DepartmentService.
func NewDepartmentService(repo DepartmentRepository) *DepartmentService {
	return &DepartmentService{
		repo: repo,
	}
}

// CreateDepartment creates a department.
func (s *DepartmentService) CreateDepartment(department *models.Department) error {
	return s.repo.Create(department)
}

// GetDepartment retrieves a department by ID.
func (s *DepartmentService) GetDepartment(id int) (*models.Department, error) {
	return s.repo.FindByID(id)
}

// GetDepartmentByCode retrieves a department by code.
func (s *DepartmentService) GetDepartmentByCode(code string) (*models.Department, error) {
	return s.repo.FindByCode(code)
}

// ListDepartments retrieves all departments.
func (s *DepartmentService) ListDepartments() ([]models.Department, error) {
	return s.repo.List()
}

// UpdateDepartment updates a department.
func (s *DepartmentService) UpdateDepartment(department *models.Department) error {
	return s.repo.Update(department)
}

// DeleteDepartment removes a department.
func (s *DepartmentService) DeleteDepartment(id int) error {
	return s.repo.Delete(id)
}
