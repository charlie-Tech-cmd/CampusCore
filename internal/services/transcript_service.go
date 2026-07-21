package services

import (
	"campuscore/internal/models"
)

// TranscriptRepository defines the persistence contract required by the service.
type TranscriptRepository interface {
	Create(*models.Transcript) error
	FindByID(int) (*models.Transcript, error)
	FindByStudent(string) ([]models.Transcript, error)
	List() ([]models.Transcript, error)
	Update(*models.Transcript) error
	Delete(int) error
}

// TranscriptService contains transcript business logic.
type TranscriptService struct {
	repo TranscriptRepository
}

// NewTranscriptService creates a TranscriptService.
func NewTranscriptService(repo TranscriptRepository) *TranscriptService {
	return &TranscriptService{
		repo: repo,
	}
}

// CreateTranscript creates a transcript.
func (s *TranscriptService) CreateTranscript(transcript *models.Transcript) error {
	return s.repo.Create(transcript)
}

// GetTranscript retrieves a transcript by ID.
func (s *TranscriptService) GetTranscript(id int) (*models.Transcript, error) {
	return s.repo.FindByID(id)
}

// GetStudentTranscripts retrieves every transcript belonging to a student.
func (s *TranscriptService) GetStudentTranscripts(studentID string) ([]models.Transcript, error) {
	return s.repo.FindByStudent(studentID)
}

// ListTranscripts retrieves all transcripts.
func (s *TranscriptService) ListTranscripts() ([]models.Transcript, error) {
	return s.repo.List()
}

// UpdateTranscript updates a transcript.
func (s *TranscriptService) UpdateTranscript(transcript *models.Transcript) error {
	return s.repo.Update(transcript)
}

// DeleteTranscript removes a transcript.
func (s *TranscriptService) DeleteTranscript(id int) error {
	return s.repo.Delete(id)
}
