package services

import "campuscore/internal/models"

// ResultRepository defines the persistence contract required by the service.
type ResultRepository interface {
	Submit(*models.Result) error
	FindByStudent(studentID string) ([]models.Result, error)
	FindByCourse(courseCode string) ([]models.Result, error)
	Update(*models.Result) error
	Delete(id int) error
}

// ResultService contains result business logic.
type ResultService struct {
	repo ResultRepository
}

// NewResultService creates a new ResultService.
func NewResultService(repo ResultRepository) *ResultService {
	return &ResultService{
		repo: repo,
	}
}

// SubmitResult submits a student's result.
func (s *ResultService) SubmitResult(result *models.Result) error {
	return s.repo.Submit(result)
}

// GetStudentResults returns all results belonging to a student.
func (s *ResultService) GetStudentResults(studentID string) ([]models.Result, error) {
	return s.repo.FindByStudent(studentID)
}

// GetCourseResults returns all results for a course.
func (s *ResultService) GetCourseResults(courseCode string) ([]models.Result, error) {
	return s.repo.FindByCourse(courseCode)
}

// UpdateResult updates an existing result.
func (s *ResultService) UpdateResult(result *models.Result) error {
	return s.repo.Update(result)
}

// DeleteResult removes a result.
func (s *ResultService) DeleteResult(id int) error {
	return s.repo.Delete(id)
}
