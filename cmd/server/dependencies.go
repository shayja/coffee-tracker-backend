// file: cmd/server/dependencies.go
package main

import (
	"coffee-tracker-backend/internal/infrastructure/auth"
	"coffee-tracker-backend/internal/infrastructure/config"
	"coffee-tracker-backend/internal/infrastructure/database"
	"coffee-tracker-backend/internal/infrastructure/http/handlers"
	"coffee-tracker-backend/internal/infrastructure/notifications"
	"coffee-tracker-backend/internal/infrastructure/repositories"
	"coffee-tracker-backend/internal/infrastructure/storage"
	"coffee-tracker-backend/internal/usecases"
	"fmt"
)

// initializeDependencies sets up all dependencies (database, repositories, use cases, handlers)
func (s *Server) initializeDependencies() error {
	// Initialize database
	db, err := database.NewSupabaseDB(s.config.DatabaseURL)
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
    var smsService notifications.SMSService

	if s.config.Env == "dev" {
		smsService = notifications.NewNoOpSMSService()
	} else {
		smsService = notifications.NewTwilioSMSService(
			"",//os.Getenv("TWILIO_ACCOUNT_SID"),
			"",//os.Getenv("TWILIO_AUTH_TOKEN"),
			"",//os.Getenv("TWILIO_FROM_NUMBER"),
		)
	 }

	// Initialize use cases
	createCoffeeUC := usecases.NewCreateCoffeeEntryUseCase(coffeeRepo)
	updateCoffeeEntryUC := usecases.NewUpdateCoffeeEntryUseCase(coffeeRepo)
	deleteCoffeeUC := usecases.NewDeleteCoffeeEntryUseCase(coffeeRepo)
	clearCoffeeEntriesUC := usecases.NewClearCoffeeEntriesUseCase(coffeeRepo)
	getCoffeeEntriesUC := usecases.NewGetCoffeeEntriesUseCase(coffeeRepo)
	getStatsUseCase := usecases.NewGetCoffeeStatsUseCase(coffeeRepo)
	getUserByIDUC := usecases.NewGetUserByIDUseCase(userRepo)
	getUserByMobileUC := usecases.NewGetUserByMobileUseCase(userRepo)
	generateOtpUC := usecases.NewGenerateOtpUseCase(authRepo, smsService, config.OtpStrength(s.config.OtpStrength))
	validateOtpUC := usecases.NewValidateOtpUseCase(authRepo, s.config.MagicOtp)
	saveRefreshTokenUC := usecases.NewSaveRefreshTokenUseCase(authRepo)
	getRefreshTokenUC := usecases.NewGetRefreshTokenUseCase(authRepo)
	deleteRefreshTokenUC := usecases.NewDeleteRefreshTokenUseCase(authRepo)

	getGenericKvUC := usecases.NewGetGenericKVUseCase(genericKvRepo)

	getProfileUC := usecases.NewGetUserProfileUseCase(userRepo)
	updateProfileUC := usecases.NewUpdateUserProfileUseCase(userRepo)
	uploadImageUC := usecases.NewUploadUserProfileImageUseCase(userRepo, storageService, s.config.ProfileImageBucket)
	deleteImageUC := usecases.NewDeleteUserProfileImageUseCase(userRepo)

	// Initialize handlers
	s.coffeeHandler = handlers.NewCoffeeEntryHandler(
		createCoffeeUC,
		getCoffeeEntriesUC,
		updateCoffeeEntryUC,
		deleteCoffeeUC,
		clearCoffeeEntriesUC,
		getStatsUseCase,
	)
	s.userSettingsHandler = handlers.NewUserSettingsHandler(
		usecases.NewGetUserSettingsUseCase(settingsRepo),
		usecases.NewUpdateUserSettingUseCase(settingsRepo),
	)
	s.healthHandler = handlers.NewHealthHandler()
	s.tokenService = auth.NewJWTService(s.config.JWTSecret, s.config.AccessTokenTTL, s.config.RefreshTokenTTL)
	s.authHandler = handlers.NewAuthHandler(
		s.tokenService,
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