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
	reportingRepo := repository.NewPostgresReportingRepository(db)
	admissionRepo := repository.NewPostgresAdmissionRepository(db)

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
	reportingService := services.NewReportingService(
		reportingRepo,
	)
	admissionService := services.NewAdmissionService(
		admissionRepo,
		notificationService,
	)

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

	reportingHandler := api.NewReportingHandler(
		reportingService,
	)

	admissionHandler := api.NewAdmissionHandler(
		admissionService,
	)

	// Prevent unused variable errors.
	_ = clearanceService

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
		"GET /reports/enrollment",
		reportingHandler.GetEnrollmentSummary,
	)

	mux.HandleFunc(
		"GET /reports/payments",
		reportingHandler.GetPaymentSummary,
	)

	mux.HandleFunc(
		"GET /reports/academic",
		reportingHandler.GetAcademicPerformanceSummary,
	)

	mux.HandleFunc(
		"GET /reports/clearance",
		reportingHandler.GetClearanceSummary,
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

	server := &http.Server{
		Addr:         ":8080",
		Handler:      middleware.Recovery(middleware.Logger(mux)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return server, worker
}
