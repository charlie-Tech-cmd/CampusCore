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
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=campuscore sslmode=disable"

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

	// Prevent unused variable errors until endpoints are added.
	_ = clearanceService

	// Handlers.
	authHandler := api.NewAuthHandler(userRepo, sessionManager)
	studentHandler := api.NewStudentHandler(
		academicService,
		ticketService,
	)
	lecturerHandler := api.NewLecturerHandler(governanceService)
	paymentHandler := api.NewPaymentHandler(paymentService)

	// Router.
	mux := http.NewServeMux()

	// Authentication.
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("/api/v1/auth/logout", authHandler.Logout)

	// Student routes.
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

	// Payment route.
	mux.Handle(
		"/api/v1/payments",
		authMiddleware.Authenticate(
			http.HandlerFunc(paymentHandler.VerifyPayment),
		),
	)

	// Lecturer routes.
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

	if err := db.Close(); err != nil {
		log.Println(err)
	}

	log.Println("Server stopped.")
}