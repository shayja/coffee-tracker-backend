package main

import (
	"log"
	"net/http"

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
	// userRepo := repositories.NewUserRepositoryImpl(db) // TODO: Implement user management

	// Initialize use cases
	createCoffeeUseCase := usecases.NewCreateCoffeeEntryUseCase(coffeeRepo)
	getCoffeeEntriesUseCase := usecases.NewGetCoffeeEntriesUseCase(coffeeRepo)
	getCoffeeStatsUseCase := usecases.NewGetCoffeeStatsUseCase(coffeeRepo)

	// Initialize handlers
	coffeeHandler := handlers.NewCoffeeEntryHandler(
		createCoffeeUseCase,
		getCoffeeEntriesUseCase,
		getCoffeeStatsUseCase,
	)
	healthHandler := handlers.NewHealthHandler()

	// Setup router
	router := mux.NewRouter()

	// Health endpoint (no auth required)
	router.HandleFunc("/health", healthHandler.Health).Methods("GET")

	// API routes with Supabase auth middleware
	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.SupabaseAuthMiddleware(cfg.SupabaseKey))
	
	// Coffee entries routes
	api.HandleFunc("/entries", coffeeHandler.CreateEntry).Methods("POST")
	api.HandleFunc("/entries", coffeeHandler.GetEntries).Methods("GET")
	api.HandleFunc("/stats", coffeeHandler.GetStats).Methods("GET")

	// Add CORS middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	log.Printf("🚀 Coffee Tracker API starting on port %s", cfg.Port)
	log.Printf("📊 Health check: http://localhost:%s/health", cfg.Port)
	log.Printf("☕ API endpoints: http://localhost:%s/api/v1/*", cfg.Port)
	
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
