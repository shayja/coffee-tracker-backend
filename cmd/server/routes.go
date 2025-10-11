package main

import (
	"net/http"
	"time"

	"coffee-tracker-backend/internal/infrastructure/http/middleware"

	"github.com/gorilla/mux"
)

// Route prefixes — keep them centralized to avoid typos and simplify maintenance
const (
	apiPrefix       = "/api/v1"
	authPrefix      = apiPrefix + "/auth"
	userPrefix      = apiPrefix + "/user"
	entriesPrefix   = apiPrefix + "/entries"
	settingsPrefix  = apiPrefix + "/settings"
	genericKVPrefix = apiPrefix + "/kv"
	statsPrefix     = apiPrefix + "/stats"
)

// setupRoutes configures all routes and their middleware
func (s *Server) setupRoutes() {
	s.router.Use(middleware.CorsMiddleware)
	s.router.Use(middleware.RequestLogger) 
	s.registerHealthRoutes()
	s.registerPublicRoutes()
	s.registerProtectedRoutes()

	s.logger.Printf("Registered routes:")
	s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
    t, _ := route.GetPathTemplate()
		s.logger.Println("  -", t)
    	return nil
	})
}

// -----------------------------
// 💚 Health Routes
// -----------------------------
func (s *Server) registerHealthRoutes() {
	s.router.HandleFunc("/health", s.healthHandler.Health).Methods(http.MethodGet)
}

// -----------------------------
// 🌍 Public (unauthenticated) Routes
// -----------------------------
func (s *Server) registerPublicRoutes() {
	api := s.router.NewRoute().Subrouter()

	api.HandleFunc(authPrefix+"/request-otp", s.authHandler.RequestOTP).Methods(http.MethodPost)
	api.HandleFunc(authPrefix+"/verify-otp", s.authHandler.VerifyOTP).Methods(http.MethodPost)
	api.HandleFunc(authPrefix+"/logout", s.authHandler.Logout).Methods(http.MethodPost)
}

// -----------------------------
// 🔒 Protected (authenticated) Routes
// -----------------------------
func (s *Server) registerProtectedRoutes() {
	api := s.router.NewRoute().Subrouter()
	api.Use(middleware.AuthMiddleware(s.tokenService))
	api.Use(middleware.UserMiddleware(s.userRepo, 5*time.Minute))

	// --- Auth routes ---
	api.HandleFunc(authPrefix+"/refresh", s.authHandler.RefreshToken).Methods(http.MethodPost)

	// --- User routes ---
	api.HandleFunc(userPrefix+"/profile", s.userHandler.GetProfile).Methods(http.MethodGet)
	api.HandleFunc(userPrefix+"/profile", s.userHandler.UpdateProfile).Methods(http.MethodPatch)
	api.HandleFunc(userPrefix+"/avatar", s.userHandler.UploadProfileImage).Methods(http.MethodPost)
	api.HandleFunc(userPrefix+"/avatar", s.userHandler.DeleteProfileImage).Methods(http.MethodDelete)

	// --- Generic KV store ---
	api.HandleFunc(genericKVPrefix, s.genericKvHandler.Get).Methods(http.MethodGet)

	// --- Coffee entries ---
	api.HandleFunc(entriesPrefix, s.coffeeHandler.GetAll).Methods(http.MethodGet)
	api.HandleFunc(entriesPrefix, s.coffeeHandler.Create).Methods(http.MethodPost)
	api.HandleFunc(entriesPrefix+"/{id}", s.coffeeHandler.Update).Methods(http.MethodPut)
	api.HandleFunc(entriesPrefix+"/{id}", s.coffeeHandler.Delete).Methods(http.MethodDelete)

	// --- Stats ---
	api.HandleFunc(statsPrefix, s.coffeeHandler.GetStats).Methods(http.MethodGet)

	// --- User settings ---
	api.HandleFunc(settingsPrefix, s.userSettingsHandler.GetAll).Methods(http.MethodGet)
	api.HandleFunc(settingsPrefix+"/{key}", s.userSettingsHandler.Update).Methods(http.MethodPatch)
}

