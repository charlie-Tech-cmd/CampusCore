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

	// Services.
	academicService := services.NewAcademicService(db)
	ticketService := services.NewTicketService(finRepo)
	clearanceService := services.NewClearanceService(finRepo)
	paymentService := services.NewPaymentService(finRepo)
	governanceService := governance.NewEngine(govRepo)

	// Prevent unused variable error until its endpoint is added.
	_ = clearanceService

	// Handlers.
	authHandler := api.NewAuthHandler(userRepo, sessionManager)

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

	// Register routes.
	mux := registerRoutes(
		authMiddleware,
		authHandler,
		studentHandler,
		lecturerHandler,
		paymentHandler,
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