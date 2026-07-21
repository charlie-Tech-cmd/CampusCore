package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"campuscore/internal/models"
)

// DepartmentManager defines the service methods required by the handler.
type DepartmentManager interface {
	CreateDepartment(*models.Department) error
	GetDepartment(int) (*models.Department, error)
	GetDepartmentByCode(string) (*models.Department, error)
	ListDepartments() ([]models.Department, error)
	UpdateDepartment(*models.Department) error
	DeleteDepartment(int) error
}

// DepartmentHandler handles department endpoints.
type DepartmentHandler struct {
	service DepartmentManager
}

// NewDepartmentHandler creates a new DepartmentHandler.
func NewDepartmentHandler(service DepartmentManager) *DepartmentHandler {
	return &DepartmentHandler{
		service: service,
	}
}

// Create handles department creation.
func (h *DepartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var department models.Department
	if err := json.NewDecoder(r.Body).Decode(&department); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateDepartment(&department); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(department)
}

// List returns all departments.
func (h *DepartmentHandler) List(w http.ResponseWriter, r *http.Request) {
	departments, err := h.service.ListDepartments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(departments)
}

// Get returns a department by ID.
func (h *DepartmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		return
	}

	department, err := h.service.GetDepartment(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(department)
}

// Update updates a department.
func (h *DepartmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var department models.Department

	if err := json.NewDecoder(r.Body).Decode(&department); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateDepartment(&department); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(department)
}

// Delete removes a department.
func (h *DepartmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteDepartment(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
