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

	// Services.
	academicService := services.NewAcademicService(db)
	ticketService := services.NewTicketService(finRepo)
	clearanceService := services.NewClearanceService(finRepo)
	paymentService := services.NewPaymentService(finRepo)
	governanceService := governance.NewEngine(govRepo)

	departmentService := services.NewDepartmentService(departmentRepo)
	facultyService := services.NewFacultyService(facultyRepo)
	registrationService := services.NewRegistrationService(
		userRepo,
		courseRepo,
		enrollmentRepo,
	)
	resultService := services.NewResultService(resultRepo)
	resultHandler := api.NewResultHandler(resultService)
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

	registrationHandler := api.NewRegistrationHandler(
		registrationService,
	)

	// Prevent unused variable errors.
	_ = clearanceService
	_ = departmentHandler
	_ = facultyHandler

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
