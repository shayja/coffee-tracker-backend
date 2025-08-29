// file: cmd/server/main.go
// Main entry point for the Coffee Tracker backend API server.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"coffee-tracker-backend/internal/infrastructure/auth"
	"coffee-tracker-backend/internal/infrastructure/config"
	"coffee-tracker-backend/internal/infrastructure/database"
	"coffee-tracker-backend/internal/infrastructure/http/handlers"
	"coffee-tracker-backend/internal/infrastructure/http/middleware"
	"coffee-tracker-backend/internal/infrastructure/repositories"
	"coffee-tracker-backend/internal/usecases"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate required environment variables
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Connect to database
	db, err := database.NewSupabaseConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	coffeeRepo := repositories.NewCoffeeEntryRepositoryImpl(db)
	userRepo := repositories.NewUserRepositoryImpl(db)
	settingsRepo := repositories.NewUserSettingsRepositoryImpl(db)
	authRepo := repositories.NewAuthRepositoryImpl(db)

	// Initialize use cases
	createCoffeeUC := usecases.NewCreateCoffeeEntryUseCase(coffeeRepo)
	editCoffeeUseCase := usecases.NewEditCoffeeEntryUseCase(coffeeRepo)
	deleteCoffeeUC := usecases.NewDeleteCoffeeEntryUseCase(coffeeRepo)
	listCoffeeUC := usecases.NewListCoffeeEntriesUseCase(coffeeRepo)
	getStatsUseCase := usecases.NewGetCoffeeStatsUseCase(coffeeRepo)


	getUserByIDUC := usecases.NewGetUserByIDUseCase(userRepo)
	getUserByMobileUC := usecases.NewGetUserByMobileUseCase(userRepo)
	genereteOtpUC := usecases.NewGenerateOtpUseCase(authRepo)
	validateOtpUC := usecases.NewValidateOtpUseCase(authRepo, cfg)
	saveRefreshTokenUC := usecases.NewSaveRefreshTokenUseCase(authRepo)
	getRefreshTokenUC := usecases.NewGetRefreshTokenUseCase(authRepo)
	deleteRefreshTokenUC := usecases.NewDeleteRefreshTokenUseCase(authRepo)

	// Initialize handlers
	coffeeHandler := handlers.NewCoffeeEntryHandler(
		createCoffeeUC,
		editCoffeeUseCase,
		deleteCoffeeUC,
		listCoffeeUC,
		getStatsUseCase,
	)

	// User settings
	getAllUC := usecases.NewGetUserSettingsUseCase(settingsRepo)
	updateUC := usecases.NewUpdateUserSettingUseCase(settingsRepo)
	userSettingsHandler := handlers.NewUserSettingsHandler(getAllUC, updateUC)

	healthHandler := handlers.NewHealthHandler()

	// authHandler
	jwtService := auth.NewJWTService(cfg.JWTSecret, 15*time.Minute, 7 * 24 * time.Hour) // 15 min access, 7 days refresh
	authHandler := handlers.NewAuthHandler(
		jwtService, 
		getUserByIDUC, 
		getUserByMobileUC,
		genereteOtpUC,
		validateOtpUC,
		saveRefreshTokenUC,
		getRefreshTokenUC,
		deleteRefreshTokenUC,
	)


	// Setup router
	router := mux.NewRouter()

	// Apply CORS middleware globally (before auth)
	router.Use(middleware.CorsMiddleware)

	// Health endpoint (no auth required)
	router.HandleFunc("/health", healthHandler.Health).Methods("GET")

	// API base router
	api := router.PathPrefix("/api/v1").Subrouter()
	// Attach logger first, so it runs before everything
	api.Use(middleware.RequestLogger)

	// ----------------- Public routes (NO auth) -----------------
	api.HandleFunc("/auth/request-otp", authHandler.RequestOTP).Methods("POST")
	api.HandleFunc("/auth/verify-otp", authHandler.VerifyOTP).Methods("POST")
	api.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")
	// ----------------- Protected routes -----------------
	protected := api.NewRoute().Subrouter()

	// User authentication middleware
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// User status middleware (e.g., check if user is active)
	protected.Use(middleware.UserMiddleware(userRepo, 5*time.Minute))

	// Auth routes (protected)
	protected.HandleFunc("/auth/token", authHandler.CreateAuthToken).Methods("GET")

	// Coffee entries routes (protected)
	protected.HandleFunc("/entries", coffeeHandler.GetEntries).Methods("GET")
	protected.HandleFunc("/entries", coffeeHandler.CreateEntry).Methods("POST")
	protected.HandleFunc("/entries/{id}", coffeeHandler.EditEntry).Methods("PUT")
	protected.HandleFunc("/entries/{id}", coffeeHandler.DeleteEntry).Methods("DELETE")
	protected.HandleFunc("/stats", coffeeHandler.GetStats).Methods("GET")


	protected.HandleFunc("/settings", userSettingsHandler.GetAll).Methods("GET")
	protected.HandleFunc("/settings/{key}", userSettingsHandler.Update).Methods("PUT")

	log.Printf("🚀 Coffee Tracker API starting on port %s", cfg.Port)

	// Fly.io specific detection
	var healthURL, apiURL string
	
	if os.Getenv("FLY_APP_NAME") != "" {
		// Production on Fly.io
		appName := os.Getenv("FLY_APP_NAME")
		healthURL = "https://" + appName + ".fly.dev/health"
		apiURL = "https://" + appName + ".fly.dev/api/v1/*"
	} else {
		// Local development
		healthURL = "http://localhost:" + cfg.Port + "/health"
		apiURL = "http://localhost:" + cfg.Port + "/api/v1/*"
	}

	log.Printf("📊 Health check: %s", healthURL)
	log.Printf("☕ API endpoints: %s", apiURL)

	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

