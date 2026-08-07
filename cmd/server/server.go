package main

import (
	"database/sql"
	"net/http"
	"time"

	"campuscore/internal/api"
	"campuscore/internal/auth"
	"campuscore/internal/governance"
	"campuscore/internal/middleware"
	"campuscore/internal/notification"
	"campuscore/internal/repository"
	"campuscore/internal/services"
)

func newServer(db *sql.DB) (*http.Server, *notification.Worker) {
	// Background worker.
	worker := notification.NewWorker(100)
	worker.Start()

	// Authentication.
	sessionManager := auth.NewSessionManager()
	authMiddleware := middleware.NewAuthGatekeeper(sessionManager)

	// Repositories.
	userRepo := repository.NewPostgresUserRepository(db)
	govRepo := repository.NewPostgresGovernanceRepository(db)
	finRepo := repository.NewPostgresFinancialRepository(db)
	departmentRepo := repository.NewPostgresDepartmentRepository(db)
	facultyRepo := repository.NewPostgresFacultyRepository(db)
	enrollmentRepo := repository.NewPostgresEnrollmentRepository(db)
	courseRepo := repository.NewPostgresCourseRepository(db)
	resultRepo := repository.NewPostgresResultRepository(db)
	attendanceRepo := repository.NewPostgresAttendanceRepository(db)
	notificationRepo := repository.NewPostgresNotificationRepository(db)
	admissionRepo := repository.NewPostgresAdmissionRepository(db)
	billingRepo := repository.NewPostgresBillingRepository(db)
	reportRepo := repository.NewPostgresReportRepository(db)

	billingService := services.NewBillingService(
		billingRepo,
	)

	notificationService := services.NewNotificationService(
		notificationRepo,
		worker,
	)

	// Services.
	academicService := services.NewAcademicService(db)
	ticketService := services.NewTicketService(finRepo)
	clearanceService := services.NewClearanceService(finRepo)
	paymentService := services.NewPaymentService(
		finRepo,
		notificationService,
	)
	governanceService := governance.NewEngine(govRepo)

	departmentService := services.NewDepartmentService(departmentRepo)
	facultyService := services.NewFacultyService(facultyRepo)
	courseService := services.NewCourseService(courseRepo)

	registrationService := services.NewRegistrationService(
		userRepo,
		courseRepo,
		enrollmentRepo,
	)
	resultService := services.NewResultService(resultRepo)

	attendanceService := services.NewAttendanceService(
		attendanceRepo,
	)
	admissionService := services.NewAdmissionService(
		admissionRepo,
		notificationService,
		billingService,
	)
	reportService := services.NewReportService(reportRepo)
	userService := services.NewUserService(userRepo)

	// Handlers.
	authHandler := api.NewAuthHandler(userRepo, sessionManager)
	refreshHandler := api.NewRefreshHandler()

	studentHandler := api.NewStudentHandler(
		academicService,
		ticketService,
	)

	lecturerHandler := api.NewLecturerHandler(
		governanceService,
	)

	paymentHandler := api.NewPaymentHandler(
		paymentService,
	)

	departmentHandler := api.NewDepartmentHandler(
		departmentService,
	)

	facultyHandler := api.NewFacultyHandler(
		facultyService,
	)

	courseHandler := api.NewCourseHandler(
		courseService,
	)

	registrationHandler := api.NewRegistrationHandler(
		registrationService,
	)

	resultHandler := api.NewResultHandler(
		resultService,
	)

	attendanceHandler := api.NewAttendanceHandler(
		attendanceService,
	)

	admissionHandler := api.NewAdmissionHandler(
		admissionService,
	)

	reportHandler := api.NewReportHandler(
		reportService)

	userHandler := api.NewUserHandler(userService)

	// Prevent unused variable errors.
	_ = clearanceService
	_ = billingService

	// Register routes.
	mux := registerRoutes(
		authMiddleware,
		authHandler,
		refreshHandler,
		studentHandler,
		lecturerHandler,
		paymentHandler,
		departmentHandler,
		facultyHandler,
		registrationHandler,
		resultHandler,
		courseHandler,
		attendanceHandler,
		admissionHandler,
	)

	mux.HandleFunc(
		"POST /admission/apply",
		admissionHandler.SubmitApplication,
	)

	mux.HandleFunc(
		"GET /admission",
		admissionHandler.ListApplications,
	)

	mux.HandleFunc(
		"GET /admission/application",
		admissionHandler.GetApplication,
	)

	mux.HandleFunc(
		"POST /admission/approve",
		admissionHandler.ApproveApplication,
	)

	mux.HandleFunc(
		"POST /admission/reject",
		admissionHandler.RejectApplication,
	)

	mux.HandleFunc(
		"/POST/reports/dashboard",
		reportHandler.DashboardSummary,
	)

	mux.HandleFunc("/students/profile", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetProfile(w, r)
		case http.MethodPut:
			userHandler.UpdateProfile(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      middleware.Recovery(middleware.Logger(mux)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return server, worker
}
