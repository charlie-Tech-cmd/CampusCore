package api

import (
	"encoding/json"
	"net/http"

	"campuscore/internal/middleware"
	"campuscore/internal/models"
)

type TranscriptManager interface {
	GenerateTranscript(string) (*models.Transcript, error)
}

type TranscriptHandler struct {
	service TranscriptManager
}

func NewTranscriptHandler(service TranscriptManager) *TranscriptHandler {
	return &TranscriptHandler{
		service: service,
	}
}

func (h *TranscriptHandler) Student(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	studentID := middleware.CurrentUserID(r)

	if studentID == "" {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	transcript, err := h.service.GenerateTranscript(studentID)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(transcript); err != nil {
		http.Error(
			w,
			"failed to encode transcript",
			http.StatusInternalServerError,
		)
		return
	}
}
