package services

import (
// "campuscore/internal/models"
)

// TranscriptService generates student transcripts.
type TranscriptService struct {
	students    StudentRepository
	courses     CourseRepository
	results     ResultRepository
	departments DepartmentRepository
	faculties   FacultyRepository
}

// NewTranscriptService creates a transcript service.
func NewTranscriptService(
	students StudentRepository,
	courses CourseRepository,
	results ResultRepository,
	departments DepartmentRepository,
	faculties FacultyRepository,
) *TranscriptService {

	return &TranscriptService{
		students:    students,
		courses:     courses,
		results:     results,
		departments: departments,
		faculties:   faculties,
	}
}
