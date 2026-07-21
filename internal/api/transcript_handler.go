package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"campuscore/internal/models"
)

// TranscriptManager defines the service contract.
type TranscriptManager interface {
	CreateTranscript(*models.Transcript) error
	GetTranscript(int) (*models.Transcript, error)
	GetStudentTranscripts(string) ([]models.Transcript, error)
	ListTranscripts() ([]models.Transcript, error)
	UpdateTranscript(*models.Transcript) error
	DeleteTranscript(int) error
}

// TranscriptHandler handles transcript endpoints.
type TranscriptHandler struct {
	service TranscriptManager
}

// NewTranscriptHandler creates a TranscriptHandler.
func NewTranscriptHandler(service TranscriptManager) *TranscriptHandler {
	return &TranscriptHandler{
		service: service,
	}
}

// Create creates a transcript.
func (h *TranscriptHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var transcript models.Transcript

	if err := json.NewDecoder(r.Body).Decode(&transcript); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateTranscript(&transcript); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transcript)
}

// Get returns a transcript by ID.
func (h *TranscriptHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid transcript id", http.StatusBadRequest)
		return
	}

	transcript, err := h.service.GetTranscript(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(transcript)
}

// List returns every transcript.
func (h *TranscriptHandler) List(w http.ResponseWriter, r *http.Request) {
	transcripts, err := h.service.ListTranscripts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transcripts)
}

// Student returns all transcripts belonging to one student.
func (h *TranscriptHandler) Student(w http.ResponseWriter, r *http.Request) {
	studentID := r.URL.Query().Get("student_id")

	transcripts, err := h.service.GetStudentTranscripts(studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transcripts)
}

// Update updates a transcript.
func (h *TranscriptHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var transcript models.Transcript

	if err := json.NewDecoder(r.Body).Decode(&transcript); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateTranscript(&transcript); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transcript)
}

// Delete deletes a transcript.
func (h *TranscriptHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid transcript id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTranscript(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
