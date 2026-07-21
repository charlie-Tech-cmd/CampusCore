package api

import (
	"campuscore/internal/models"
	// "encoding/json"
	// "net/http"
	// "strconv"
)

type TimetableService interface {
	Create(entry *models.Timetable) error
	Get(id int) (*models.Timetable, error)
	List() ([]models.Timetable, error)
	ListByDepartment(departmentID int) ([]models.Timetable, error)
	ListByLecturer(lecturerID string) ([]models.Timetable, error)
	ListByCourse(courseCode string) ([]models.Timetable, error)
	ListByLevel(level int) ([]models.Timetable, error)
	Update(entry *models.Timetable) error
	Delete(id int) error
}

// Timetable represents a scheduled class for a course.
type TimetableHandler struct {
	service TimetableService
}

func NewTimetableHandler(service TimetableService) *TimetableHandler {
	return &TimetableHandler{
		service: service,
	}
}
