package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"coffee-tracker-backend/internal/domain/repositories"
	"coffee-tracker-backend/internal/infrastructure/auth"
	"coffee-tracker-backend/internal/infrastructure/config"
	"coffee-tracker-backend/internal/infrastructure/http/handlers"

	"github.com/gorilla/mux"
)

// Server encapsulates the HTTP server and its dependencies
type Server struct {
	config            	*config.Config
	router            	*mux.Router
	httpServer        	*http.Server
	logger            	*log.Logger
	userHandler 	  	*handlers.UserHandler
	genericKvHandler	*handlers.GenericKVHandler
	coffeeHandler     	*handlers.CoffeeEntryHandler
	userSettingsHandler *handlers.UserSettingsHandler
	healthHandler     	*handlers.HealthHandler
	authHandler       	*handlers.AuthHandler
	taperingHandler		*handlers.TaperingJourneyHandler
	jwtService        	*auth.JWTService
	userRepo          	repositories.UserRepository
	//db                *database.Supabase // Added to manage DB connection
}

// NewServer initializes a new Server instance with all dependencies
func NewServer() (*Server, error) {
	// Initialize logger
	logger := log.New(os.Stdout, "coffee-tracker: ", log.LstdFlags|log.Lshortfile)

	// Load configuration
	cfg := config.Load()
	// Validate required environment variables
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Initialize server
	server := &Server{
		config: cfg,
		router: mux.NewRouter(),
		logger: logger,
	}

	// Initialize dependencies and routes
	if err := server.initializeDependencies(); err != nil {
		return nil, err
	}
	server.setupRoutes()

	// Configure HTTP server
	server.httpServer = &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      server.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return server, nil
}

// Start runs the HTTP server with graceful shutdown
func (s *Server) Start() error {
	// Log server information
	s.logServerInfo()

	// Create a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		s.logger.Printf("Starting server on port %s", s.config.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	s.logger.Println("Received shutdown signal. Initiating graceful shutdown...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Close database connection
	// if err := s.db.Close(); err != nil {
	// 	s.logger.Printf("Error closing database connection: %v", err)
	// }

	// Attempt graceful shutdown
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Printf("Server shutdown error: %v", err)
		return err
	}

	s.logger.Println("Server gracefully stopped")
	return nil
}

// logServerInfo logs the server endpoints based on environment
func (s *Server) logServerInfo() {
	var healthURL, apiURL string
	if os.Getenv("FLY_APP_NAME") != "" {
		appName := os.Getenv("FLY_APP_NAME")
		healthURL = "https://" + appName + ".fly.dev/health"
		apiURL = "https://" + appName + ".fly.dev/api/v1/*"
	} else {
		healthURL = "http://localhost:" + s.config.Port + "/health"
		apiURL = "http://localhost:" + s.config.Port + "/api/v1/*"
	}
	s.logger.Printf("🚀 Coffee Tracker API starting...")
	s.logger.Printf("📊 Health check: %s", healthURL)
	s.logger.Printf("☕ API endpoints: %s", apiURL)
}