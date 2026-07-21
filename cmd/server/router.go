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
	refreshHandler *api.RefreshHandler,
	studentHandler *api.StudentHandler,
	lecturerHandler *api.LecturerHandler,
	paymentHandler *api.PaymentHandler,
	departmentHandler *api.DepartmentHandler,
	facultyHandler *api.FacultyHandler,
	registrationHandler *api.RegistrationHandler,
	resultHandler *api.ResultHandler,
) *http.ServeMux {

	mux := http.NewServeMux()

	// Results
	mux.Handle(
		"/api/v1/results/submit",
		authMiddleware.Authenticate(
			http.HandlerFunc(resultHandler.Submit),
		),
	)

	mux.Handle(
		"/api/v1/results/student",
		authMiddleware.Authenticate(
			http.HandlerFunc(resultHandler.StudentResults),
		),
	)

	mux.Handle(
		"/api/v1/results/course",
		authMiddleware.Authenticate(
			http.HandlerFunc(resultHandler.CourseResults),
		),
	)

	mux.Handle(
		"/api/v1/results/update",
		authMiddleware.Authenticate(
			http.HandlerFunc(resultHandler.Update),
		),
	)

	mux.Handle(
		"/api/v1/results/delete",
		authMiddleware.Authenticate(
			http.HandlerFunc(resultHandler.Delete),
		),
	)

	// Home
	mux.HandleFunc("/", homeHandler)

	// Health
	mux.HandleFunc("/health", healthHandler)

	// Authentication
	mux.HandleFunc("/api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("/api/v1/auth/logout", authHandler.Logout)
	mux.HandleFunc("/api/v1/auth/refresh", refreshHandler.RefreshToken)

	mux.Handle(
		"/api/v1/auth/me",
		authMiddleware.Authenticate(
			http.HandlerFunc(authHandler.Me),
		),
	)

	// Faculties
	mux.Handle(
		"/api/v1/faculties",
		authMiddleware.Authenticate(
			http.HandlerFunc(facultyHandler.List),
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

	mux.Handle(
		"/api/v1/student/profile",
		authMiddleware.Authenticate(
			authMiddleware.RequireRole("student")(
				http.HandlerFunc(studentHandler.GetProfile),
			),
		),
	)

	mux.Handle(
		"/api/v1/student/profile/update",
		authMiddleware.Authenticate(
			authMiddleware.RequireRole("student")(
				http.HandlerFunc(studentHandler.UpdateProfile),
			),
		),
	)

	mux.HandleFunc(
		"/students/register-course",
		registrationHandler.RegisterCourse,
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

	// Departments
	mux.Handle(
		"/api/v1/departments",
		authMiddleware.Authenticate(
			http.HandlerFunc(departmentHandler.List),
		),
	)

	mux.Handle(
		"/api/v1/departments/create",
		authMiddleware.Authenticate(
			http.HandlerFunc(departmentHandler.Create),
		),
	)

	mux.Handle(
		"/api/v1/departments/get",
		authMiddleware.Authenticate(
			http.HandlerFunc(departmentHandler.Get),
		),
	)

	mux.Handle(
		"/api/v1/departments/update",
		authMiddleware.Authenticate(
			http.HandlerFunc(departmentHandler.Update),
		),
	)

	mux.Handle(
		"/api/v1/departments/delete",
		authMiddleware.Authenticate(
			http.HandlerFunc(departmentHandler.Delete),
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
