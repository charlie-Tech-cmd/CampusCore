package main

import (
	"context"
	"database/sql"
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

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting CampusCore...")

	// Connect to PostgreSQL.
	connStr := "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=campuscore sslmode=disable"

	log.Println(connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(15 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected.")

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

	// Register all routes.
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println(err)
	}

	if err := db.Close(); err != nil {
		log.Println(err)
	}

	log.Println("Server stopped.")
}