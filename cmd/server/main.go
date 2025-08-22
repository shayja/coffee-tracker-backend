package main

import (
	"log"
	"net/http"
	"time"

	"coffee-tracker-backend/internal/infrastructure/config"
	"coffee-tracker-backend/internal/infrastructure/database"
	"coffee-tracker-backend/internal/infrastructure/http/handlers"
	"coffee-tracker-backend/internal/infrastructure/http/middleware"
	"coffee-tracker-backend/internal/infrastructure/repositories"
	"coffee-tracker-backend/internal/services"
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

	// Initialize use cases
	createCoffeeUseCase := usecases.NewCreateCoffeeEntryUseCase(coffeeRepo)
	editCoffeeUseCase := usecases.NewEditCoffeeEntryUseCase(coffeeRepo)
	deleteCoffeeUseCase := usecases.NewDeleteCoffeeEntryUseCase(coffeeRepo)
	getCoffeeEntriesUseCase := usecases.NewGetCoffeeEntriesUseCase(coffeeRepo)
	getCoffeeStatsUseCase := usecases.NewGetCoffeeStatsUseCase(coffeeRepo)

	// Initialize handlers
	coffeeHandler := handlers.NewCoffeeEntryHandler(
		createCoffeeUseCase,
		editCoffeeUseCase,
		deleteCoffeeUseCase,
		getCoffeeEntriesUseCase,
		getCoffeeStatsUseCase,
	)
	healthHandler := handlers.NewHealthHandler()
	// Initialize auth service
	authService := services.NewAuthService(repositories.NewAuthRepositoryImpl(db))
	userService := services.NewUserService(userRepo)
	authHandler := handlers.NewAuthHandler(cfg.JWTSecret, authService, userService)


	// Setup router
	router := mux.NewRouter()

	// Apply CORS middleware globally (before auth)
	router.Use(corsMiddleware)

	// Health endpoint (no auth required)
	router.HandleFunc("/health", healthHandler.Health).Methods("GET")

	// API base router
	api := router.PathPrefix("/api/v1").Subrouter()
	// Attach logger first, so it runs before everything
	api.Use(middleware.RequestLogger)

	// ----------------- Public routes (NO auth) -----------------
	api.HandleFunc("/auth/request-otp", authHandler.RequestOTP).Methods("POST")
	api.HandleFunc("/auth/verify-otp", authHandler.VerifyOTP).Methods("POST")

	// ----------------- Protected routes -----------------
	protected := api.NewRoute().Subrouter()

	// User authentication middleware
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// User status middleware (e.g., check if user is active)
	protected.Use(middleware.UserMiddleware(userRepo, 5*time.Minute))

	// Auth routes (protected)
	protected.HandleFunc("/auth/token", authHandler.CreateAuthToken).Methods("GET")
	protected.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")

	// Coffee entries routes (protected)
	protected.HandleFunc("/entries", coffeeHandler.GetEntries).Methods("GET")
	protected.HandleFunc("/entries", coffeeHandler.CreateEntry).Methods("POST")
	protected.HandleFunc("/entries/{id}", coffeeHandler.EditEntry).Methods("PUT")
	protected.HandleFunc("/entries/{id}", coffeeHandler.DeleteEntry).Methods("DELETE")
	protected.HandleFunc("/stats", coffeeHandler.GetStats).Methods("GET")

	// Dev-only: print JWT
	// if cfg.Env == "dev" {
	// 	printJWT(cfg.JWTSecret, uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"))
	// }

	log.Printf("🚀 Coffee Tracker API starting on port %s", cfg.Port)
	log.Printf("📊 Health check: http://localhost:%s/health", cfg.Port)
	log.Printf("☕ API endpoints: http://localhost:%s/api/v1/*", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// corsMiddleware allows cross-origin requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// printJWT generates and prints a JWT token for dev testing
/*
func printJWT(secret string, userID uuid.UUID) {
	token, err := auth.GenerateJWT(secret, userID)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Dev JWT token:", token)
}*/
