package main

import (
	"coffee-tracker-backend/internal/infrastructure/auth"
	"coffee-tracker-backend/internal/infrastructure/database"
	"coffee-tracker-backend/internal/infrastructure/http/handlers"
	"coffee-tracker-backend/internal/infrastructure/repositories"
	"coffee-tracker-backend/internal/infrastructure/storage"
	"coffee-tracker-backend/internal/usecases"
	"fmt"
	"time"
)

// initializeDependencies sets up all dependencies (database, repositories, use cases, handlers)
	func (s *Server) initializeDependencies() error {
		// Initialize database
		db, err := database.NewSupabaseConnection(s.config.DatabaseURL)
		if err != nil {
			return err
	}

	// Initialize repositories
	coffeeRepo := repositories.NewCoffeeEntryRepositoryImpl(db)
	userRepo := repositories.NewUserRepositoryImpl(db)
	settingsRepo := repositories.NewUserSettingsRepositoryImpl(db)
	authRepo := repositories.NewAuthRepositoryImpl(db)
	genericKvRepo := repositories.NewGenericKVRepositoryImpl(db)

	// Initialize Supabase Storage client
	if s.config.StorageURL == "" || s.config.ServiceRoleKey == "" {
		s.logger.Fatal(fmt.Errorf("invalid storage configuration: URL or API key missing"))
	}

    storageService := storage.NewSupabaseStorageService(s.config.StorageURL, s.config.ServiceRoleKey)
    

	// Initialize use cases
	createCoffeeUC := usecases.NewCreateCoffeeEntryUseCase(coffeeRepo)
	editCoffeeUseCase := usecases.NewEditCoffeeEntryUseCase(coffeeRepo)
	deleteCoffeeUC := usecases.NewDeleteCoffeeEntryUseCase(coffeeRepo)
	listCoffeeUC := usecases.NewListCoffeeEntriesUseCase(coffeeRepo)
	getStatsUseCase := usecases.NewGetCoffeeStatsUseCase(coffeeRepo)
	getUserByIDUC := usecases.NewGetUserByIDUseCase(userRepo)
	getUserByMobileUC := usecases.NewGetUserByMobileUseCase(userRepo)
	generateOtpUC := usecases.NewGenerateOtpUseCase(authRepo)
	validateOtpUC := usecases.NewValidateOtpUseCase(authRepo, s.config)
	saveRefreshTokenUC := usecases.NewSaveRefreshTokenUseCase(authRepo)
	getRefreshTokenUC := usecases.NewGetRefreshTokenUseCase(authRepo)
	deleteRefreshTokenUC := usecases.NewDeleteRefreshTokenUseCase(authRepo)

	getGenericKvUC := usecases.NewGetGenericKVUseCase(genericKvRepo)

	getProfileUC := usecases.NewGetUserProfileUseCase(userRepo)
	updateProfileUC := usecases.NewUpdateUserProfileUseCase(userRepo)
	uploadImageUC := usecases.NewUploadUserProfileImageUseCase(userRepo, storageService, s.config)
	deleteImageUC := usecases.NewDeleteUserProfileImageUseCase(userRepo)

	// Initialize handlers
	s.coffeeHandler = handlers.NewCoffeeEntryHandler(
		createCoffeeUC,
		editCoffeeUseCase,
		deleteCoffeeUC,
		listCoffeeUC,
		getStatsUseCase,
	)
	s.userSettingsHandler = handlers.NewUserSettingsHandler(
		usecases.NewGetUserSettingsUseCase(settingsRepo),
		usecases.NewUpdateUserSettingUseCase(settingsRepo),
	)
	s.healthHandler = handlers.NewHealthHandler()
	s.jwtService = auth.NewJWTService(s.config.JWTSecret, 15*time.Minute, 7*24*time.Hour)
	s.authHandler = handlers.NewAuthHandler(
		s.jwtService,
		getUserByIDUC,
		getUserByMobileUC,
		generateOtpUC,
		validateOtpUC,
		saveRefreshTokenUC,
		getRefreshTokenUC,
		deleteRefreshTokenUC,
	)
	s.userRepo = userRepo

	s.genericKvHandler = handlers.NewGenericKVHandler(getGenericKvUC)

	s.userHandler = handlers.NewUserHandler(
		getProfileUC,
		updateProfileUC,
		uploadImageUC, 
		deleteImageUC,
	)

	return nil
}