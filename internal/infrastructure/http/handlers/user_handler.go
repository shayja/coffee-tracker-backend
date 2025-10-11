// file: internal/infrastructure/http/handlers/user_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	httpUtils "coffee-tracker-backend/internal/infrastructure/http"
	"coffee-tracker-backend/internal/infrastructure/http/models"
	"coffee-tracker-backend/internal/usecases"
)

type UserHandler struct {
	getProfileUC    *usecases.GetUserProfileUseCase
	updateProfileUC *usecases.UpdateUserProfileUseCase
	uploadImageUC   *usecases.UploadUserProfileImageUseCase
	deleteImageUC   *usecases.DeleteUserProfileImageUseCase
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
	httpUtils.LogRequest(r)

	userID, ok := httpUtils.GetUserIDOrAbort(w, r)
	if !ok {
		return
	}

	profile, err := h.getProfileUC.Execute(r.Context(), userID)
	if err != nil {
		httpUtils.WriteError(w, http.StatusInternalServerError, "Failed to load profile", err.Error())
		return
	}

	httpUtils.WriteJSON(w, http.StatusOK, profile)
}

// PATCH /profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	httpUtils.LogRequest(r)

	userID, ok := httpUtils.GetUserIDOrAbort(w, r)
	if !ok {
		return
	}

	var req models.UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpUtils.WriteError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	err := h.updateProfileUC.Execute(r.Context(), userID, &req)
	if err != nil {
		switch err {
		case usecases.ErrInvalidInput:
			httpUtils.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			httpUtils.WriteError(w, http.StatusInternalServerError, "Failed to update profile", err.Error())
		}
		return
	}

	httpUtils.WriteJSON(w, http.StatusOK, req)
}

// POST /profile/image
func (h *UserHandler) UploadProfileImage(w http.ResponseWriter, r *http.Request) {
	httpUtils.LogRequest(r)

	userID, ok := httpUtils.GetUserIDOrAbort(w, r)
	if !ok {
		return
	}

	err := r.ParseMultipartForm(5 << 20) // 5 MB
	if err != nil {
		httpUtils.WriteError(w, http.StatusBadRequest, "Invalid or too large file", err.Error())
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		httpUtils.WriteError(w, http.StatusBadRequest, "Failed to read file", err.Error())
		return
	}
	defer file.Close()

	filename := header.Filename
	avatarURL, err := h.uploadImageUC.Execute(r.Context(), userID, filename, file)
	if err != nil {
		httpUtils.WriteError(w, http.StatusInternalServerError, "Failed to upload file", err.Error())
		return
	}

	httpUtils.WriteJSON(w, http.StatusOK, map[string]string{"avatar_url": avatarURL})
}

// DELETE /profile/image
func (h *UserHandler) DeleteProfileImage(w http.ResponseWriter, r *http.Request) {
	httpUtils.LogRequest(r)

	userID, ok := httpUtils.GetUserIDOrAbort(w, r)
	if !ok {
		return
	}

	err := h.deleteImageUC.Execute(r.Context(), userID)
	if err != nil {
		switch err {
		case usecases.ErrNotFound:
			httpUtils.WriteError(w, http.StatusNotFound, err.Error())
		default:
			httpUtils.WriteError(w, http.StatusInternalServerError, "Failed to delete image", err.Error())
		}
		return
	}

	httpUtils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Profile image deleted"})
}
