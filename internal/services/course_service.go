package services

import "campuscore/internal/models"

// CourseRepository defines the persistence contract required by the service.
type CourseServiceRepository interface {
	Create(*models.Course) error
	FindByCode(string) (*models.Course, error)
	GetAll() ([]models.Course, error)
	GetByDepartment(int) ([]models.Course, error)
	Update(*models.Course) error
	Delete(string) error
}

// CourseService contains course business logic.
type CourseService struct {
	repo CourseServiceRepository
}

// NewCourseService creates a CourseService.
func NewCourseService(repo CourseServiceRepository) *CourseService {
	return &CourseService{
		repo: repo,
	}
}

// CreateCourse creates a new course.
func (s *CourseService) CreateCourse(course *models.Course) error {
	return s.repo.Create(course)
}

// GetCourse retrieves a course by code.
func (s *CourseService) GetCourse(code string) (*models.Course, error) {
	return s.repo.FindByCode(code)
}

// ListCourses retrieves every course.
func (s *CourseService) ListCourses() ([]models.Course, error) {
	return s.repo.GetAll()
}

// ListDepartmentCourses retrieves courses belonging to a department.
func (s *CourseService) ListDepartmentCourses(departmentID int) ([]models.Course, error) {
	return s.repo.GetByDepartment(departmentID)
}

// UpdateCourse updates a course.
func (s *CourseService) UpdateCourse(course *models.Course) error {
	return s.repo.Update(course)
}

// DeleteCourse removes a course.
func (s *CourseService) DeleteCourse(code string) error {
	return s.repo.Delete(code)
}
