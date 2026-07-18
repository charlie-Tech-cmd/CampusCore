package main

import (
	"encoding/json"
	"net/http"

	"campuscore/internal/api"
	"campuscore/internal/middleware"
)

func registerRoutes(
	authMiddleware *middleware.AuthGatekeeper,
	authHandler *api.AuthHandler,
	studentHandler *api.StudentHandler,
	lecturerHandler *api.LecturerHandler,
	paymentHandler *api.PaymentHandler,
) *http.ServeMux {

	mux := http.NewServeMux()

	// Home
	mux.HandleFunc("/", homeHandler)

	// Health
	mux.HandleFunc("/health", healthHandler)

	// Authentication
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("/api/v1/auth/logout", authHandler.Logout)

	mux.Handle(
		"/api/v1/auth/me",
		authMiddleware.Authenticate(
			http.HandlerFunc(authHandler.Me),
		),
	)
	// Student
	mux.Handle(
		"/api/v1/student/courses/register",
		authMiddleware.Authenticate(
			authMiddleware.RequireRole("student")(
				http.HandlerFunc(studentHandler.RegisterCourse),
			),
		),
	)

	mux.Handle(
		"/api/v1/student/support/tickets",
		authMiddleware.Authenticate(
			authMiddleware.RequireRole("student")(
				http.HandlerFunc(studentHandler.SubmitTicket),
			),
		),
	)

	// Payments
	mux.Handle(
		"/api/v1/payments",
		authMiddleware.Authenticate(
			http.HandlerFunc(paymentHandler.VerifyPayment),
		),
	)

	// Lecturer
	mux.Handle(
		"/api/v1/faculty/results/advance",
		authMiddleware.Authenticate(
			authMiddleware.RequireRole(
				"lecturer",
				"HOD",
				"dean",
				"admin",
			)(
				http.HandlerFunc(lecturerHandler.AdvanceApproval),
			),
		),
	)

	mux.Handle(
		"/api/v1/faculty/results/reject",
		authMiddleware.Authenticate(
			authMiddleware.RequireRole(
				"lecturer",
				"HOD",
				"dean",
				"admin",
			)(
				http.HandlerFunc(lecturerHandler.RejectApproval),
			),
		),
	)

	return mux
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]any{
		"application": "CampusCore API",
		"version":     "v1",
		"status":      "running",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status": "healthy",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
