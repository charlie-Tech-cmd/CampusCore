package services

import "campuscore/internal/models"

// TimetableRepository defines timetable persistence operations.
type TimetableRepository interface {
	Create(entry *models.Timetable) error
	FindByID(id int) (*models.Timetable, error)
	GetAll() ([]models.Timetable, error)
	GetByDepartment(departmentID int) ([]models.Timetable, error)
	GetByLecturer(lecturerID string) ([]models.Timetable, error)
	GetByCourse(courseCode string) ([]models.Timetable, error)
	GetByLevel(level int) ([]models.Timetable, error)
	Update(entry *models.Timetable) error
	Delete(id int) error
}

// TimetableService coordinates timetable operations.
type TimetableService struct {
	repo TimetableRepository
}

// NewTimetableService creates a timetable service.
func NewTimetableService(repo TimetableRepository) *TimetableService {
	return &TimetableService{
		repo: repo,
	}
}

// Create creates a timetable entry.
func (s *TimetableService) Create(entry *models.Timetable) error {
	return s.repo.Create(entry)
}

// Get returns a timetable entry by ID.
func (s *TimetableService) Get(id int) (*models.Timetable, error) {
	return s.repo.FindByID(id)
}

// List returns all timetable entries.
func (s *TimetableService) List() ([]models.Timetable, error) {
	return s.repo.GetAll()
}

// ListByDepartment returns a department timetable.
func (s *TimetableService) ListByDepartment(departmentID int) ([]models.Timetable, error) {
	return s.repo.GetByDepartment(departmentID)
}

// ListByLecturer returns a lecturer timetable.
func (s *TimetableService) ListByLecturer(lecturerID string) ([]models.Timetable, error) {
	return s.repo.GetByLecturer(lecturerID)
}

// ListByCourse returns a course timetable.
func (s *TimetableService) ListByCourse(courseCode string) ([]models.Timetable, error) {
	return s.repo.GetByCourse(courseCode)
}

// ListByLevel returns a level timetable.
func (s *TimetableService) ListByLevel(level int) ([]models.Timetable, error) {
	return s.repo.GetByLevel(level)
}

// Update modifies a timetable entry.
func (s *TimetableService) Update(entry *models.Timetable) error {
	return s.repo.Update(entry)
}

// Delete removes a timetable entry.
func (s *TimetableService) Delete(id int) error {
	return s.repo.Delete(id)
}
