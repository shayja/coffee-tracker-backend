// file: cmd/server/server.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"coffee-tracker-backend/internal/infrastructure/auth"
	"coffee-tracker-backend/internal/infrastructure/config"
	"coffee-tracker-backend/internal/infrastructure/http/handlers"
	"coffee-tracker-backend/internal/repositories"

	"github.com/gorilla/mux"
)

// Server encapsulates the HTTP server and its dependencies
type Server struct {
	config              *config.Config
	router              *mux.Router
	httpServer          *http.Server
	logger              *log.Logger
	userHandler         *handlers.UserHandler
	genericKvHandler    *handlers.GenericKVHandler
	coffeeHandler       *handlers.CoffeeEntryHandler
	userSettingsHandler *handlers.UserSettingsHandler
	healthHandler       *handlers.HealthHandler
	authHandler         *handlers.AuthHandler
	tokenService        auth.TokenService
	userRepo            repositories.UserRepository
}

// NewServer initializes a new Server instance with all dependencies
func NewServer() (*Server, error) {
	// Initialize logger
	logger := log.New(os.Stdout, "coffee-tracker: ", log.LstdFlags|log.Lshortfile)

	// Load configuration
	cfg, err := config.Load() // already validated
	if err != nil {
		logger.Fatalf("❌ Invalid configuration: %v", err)
	}

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

// Start runs the HTTP server (blocking)
func (s *Server) Start() error {
	s.logServerInfo()
	s.logger.Printf("🚀 Starting server on port %s", s.config.Port)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Println("🧹 Shutting down server...")
	// Attempt graceful shutdown
	return s.httpServer.Shutdown(ctx)
}

// logServerInfo logs environment-based URLs for debugging
func (s *Server) logServerInfo() {
	var baseURL string
	if appName := os.Getenv("FLY_APP_NAME"); appName != "" {
		baseURL = "https://" + appName + ".fly.dev"
	} else {
		baseURL = "http://localhost:" + s.config.Port
	}
	s.logger.Printf("🌍 Environment: %s", s.config.Env)
	s.logger.Printf("📊 Health check: %s/health", baseURL)
	s.logger.Printf("☕ API base: %s/api/v1/*", baseURL)
}
