// file: cmd/server/routes.go
package main

import (
	"net/http"
	"time"

	"coffee-tracker-backend/internal/infrastructure/http/middleware"
)

// setupRoutes configures the API routes with middleware
func (s *Server) setupRoutes() {
	// Apply global CORS middleware
	s.router.Use(middleware.CorsMiddleware)

	// Health endpoint (no auth)
	s.router.HandleFunc("/health", s.healthHandler.Health).Methods(http.MethodGet)

	// API v1 router
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.RequestLogger)

	// Public routes
	api.HandleFunc("/auth/request-otp", s.authHandler.RequestOTP).Methods(http.MethodPost)
	api.HandleFunc("/auth/verify-otp", s.authHandler.VerifyOTP).Methods(http.MethodPost)
	api.HandleFunc("/auth/refresh", s.authHandler.RefreshToken).Methods(http.MethodPost)

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(s.config.JWTSecret))
	protected.Use(middleware.UserMiddleware(s.userRepo, 5*time.Minute))

	// Auth routes
	protected.HandleFunc("/auth/token", s.authHandler.CreateAuthToken).Methods(http.MethodGet)

	// User profile routes
	protected.HandleFunc("/user/profile", s.userHandler.GetProfile).Methods(http.MethodGet)
	protected.HandleFunc("/user/profile", s.userHandler.UpdateProfile).Methods(http.MethodPatch)

	// User avatar routes
	protected.HandleFunc("/user/avatar", s.userHandler.UploadProfileImage).Methods(http.MethodPost)
	protected.HandleFunc("/user/avatar", s.userHandler.DeleteProfileImage).Methods(http.MethodDelete)

	protected.HandleFunc("/kv", s.genericKvHandler.Get).Methods(http.MethodGet)


	// Coffee entry routes
	protected.HandleFunc("/entries", s.coffeeHandler.GetEntries).Methods(http.MethodGet)
	protected.HandleFunc("/entries", s.coffeeHandler.CreateEntry).Methods(http.MethodPost)
	protected.HandleFunc("/entries/{id}", s.coffeeHandler.EditEntry).Methods(http.MethodPut)
	protected.HandleFunc("/entries/{id}", s.coffeeHandler.DeleteEntry).Methods(http.MethodDelete)
	protected.HandleFunc("/stats", s.coffeeHandler.GetStats).Methods(http.MethodGet)

	// User settings routes
	protected.HandleFunc("/settings", s.userSettingsHandler.GetAll).Methods(http.MethodGet)
	protected.HandleFunc("/settings/{key}", s.userSettingsHandler.Update).Methods(http.MethodPatch)
}
