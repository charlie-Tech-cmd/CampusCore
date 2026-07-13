package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"campuscore/internal/api"
	"campuscore/internal/auth"
	"campuscore/internal/governance"
	"campuscore/internal/middleware"
	"campuscore/internal/notification"
	"campuscore/internal/repository"
	"campuscore/internal/services"
)

func main() {
	log.Println("Starting CampusCore...")

	db := mustConnectDB()
	defer db.Close()

	// Start background worker.
	worker := notification.NewWorker(100)
	worker.Start()
	defer worker.Stop(context.Background())

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

	// Prevent unused variable error until endpoint is added.
	_ = clearanceService

	// Handlers.
	authHandler := api.NewAuthHandler(userRepo, sessionManager)
	studentHandler := api.NewStudentHandler(
		academicService,
		ticketService,
	)
	lecturerHandler := api.NewLecturerHandler(governanceService)
	paymentHandler := api.NewPaymentHandler(paymentService)

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

	go func() {
		log.Println("Server listening on http://localhost:8080")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println(err)
	}

	log.Println("Server stopped.")
}