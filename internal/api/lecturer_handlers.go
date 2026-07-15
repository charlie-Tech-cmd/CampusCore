package api

import (
	"encoding/json"
	"net/http"
	
	"campuscore/internal/auth"
	// "campuscore/internal/governance"
	"campuscore/internal/middleware"
	"campuscore/internal/models"
)

// LecturerHandler coordinates network delivery for grade processing and workflow changes
type LecturerHandler struct {
	govEngine GovernanceEngine
}

// NewLecturerHandler instantiates our lecturer endpoint controller
func NewLecturerHandler(ge GovernanceEngine) *LecturerHandler {
	return &LecturerHandler{
		govEngine: ge,
	}
}

// WorkflowAdvanceRequest defines the expected JSON data layout to step a batch forward
type WorkflowAdvanceRequest struct {
	CourseCode string `json:"course_code"`
}

// WorkflowRejectRequest defines the expected JSON payload to cascade a grade sheet backward
type WorkflowRejectRequest struct {
	CourseCode string `json:"course_code"`
	Remarks    string `json:"remarks"` // Mandatory explanation text required to audit rejections
}

// AdvanceApproval processes structural workflow staging upgrades
func (h *LecturerHandler) AdvanceApproval(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"error": "Method not allowed. Use POST."}`))
		return
	}

	// 1. Recover the active user session context details
	sessionVal := r.Context().Value(middleware.UserContextKey)
	if sessionVal == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "Unidentified administrative context environment."}`))
		return
	}
	activeSession := sessionVal.(*auth.Session)

	// 2. Decode the execution parameters from the body
	var req WorkflowAdvanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "Malformed request payload parameters."}`))
		return
	}

	// 3. Delegate state matching transitions down to our core governance engine
	err := h.govEngine.ProcessApprovalAdvance(req.CourseCode, models.UserRole(activeSession.Role), activeSession.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Academic record workflow advanced successfully."}`))
}

// RejectApproval processes structural rollback loops down the governance tiers
func (h *LecturerHandler) RejectApproval(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"error": "Method not allowed. Use POST."}`))
		return
	}

	// 1. Recover session values to identify who is triggering the audit rejection
	sessionVal := r.Context().Value(middleware.UserContextKey)
	if sessionVal == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "Unidentified administrative context environment."}`))
		return
	}
	activeSession := sessionVal.(*auth.Session)

	// 2. Decode the rejection parameters and the mandatory audit reason string
	var req WorkflowRejectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "Malformed request payload parameters."}`))
		return
	}

	// 3. Fire the rejection rule sequence inside our state engine layer
	err := h.govEngine.ProcessApprovalRejection(req.CourseCode, models.UserRole(activeSession.Role), activeSession.UserID, req.Remarks)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Academic record batch successfully rolled back to previous tier. Author notified."}`))
}