// file: internal/infrastructure/http/handlers/user_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	http_utils "coffee-tracker-backend/internal/infrastructure/http"
	"coffee-tracker-backend/internal/infrastructure/http/dto"
	"coffee-tracker-backend/internal/usecases"
)

type UserHandler struct {
	getProfileUC      *usecases.GetUserProfileUseCase
	updateProfileUC   *usecases.UpdateUserProfileUseCase
	uploadImageUC     *usecases.UploadUserProfileImageUseCase
	deleteImageUC     *usecases.DeleteUserProfileImageUseCase
}

func NewUserHandler(
	getProfileUC *usecases.GetUserProfileUseCase,
	updateProfileUC *usecases.UpdateUserProfileUseCase,
	uploadImageUC *usecases.UploadUserProfileImageUseCase,
	deleteImageUC *usecases.DeleteUserProfileImageUseCase,
) *UserHandler {
	return &UserHandler{
		getProfileUC:    getProfileUC,
		updateProfileUC: updateProfileUC,
		uploadImageUC:   uploadImageUC,
		deleteImageUC:   deleteImageUC,
	}
}

// GET /profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok { return }

	profile, err := h.getProfileUC.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// PATCH /profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok { return }

	var req dto.UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.updateProfileUC.Execute(r.Context(), userID, &req)
	if err != nil {
		switch err {
		case usecases.ErrInvalidInput:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

// POST /profile/image
func (h *UserHandler) UploadProfileImage(w http.ResponseWriter, r *http.Request) {
    userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok { return }

    // Parse the multipart form
    err := r.ParseMultipartForm(5 << 20) // 5 MB max
    if err != nil {
        http.Error(w, "File too big or invalid", http.StatusBadRequest)
        return
    }

    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Failed to read file", http.StatusBadRequest)
        return
    }
    defer file.Close()

	// Pass to service the original filename, the service will Generate a random filename with same extension.
	filename := header.Filename

    // Call use case to store in Supabase (or any storage)
    avatarURL, err := h.uploadImageUC.Execute(r.Context(), userID, filename, file)
    if err != nil {
        http.Error(w, "Failed to upload file", http.StatusInternalServerError)
        return
    }

    // Return the new URL
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"avatar_url": avatarURL})
}

// DELETE /profile/image
func (h *UserHandler) DeleteProfileImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok { return }

	err := h.deleteImageUC.Execute(r.Context(), userID)
	if err != nil {
		switch err {
		case usecases.ErrNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
