// file: internal/infrastructure/repositories/user_repository_impl.go
package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/infrastructure/http/models"
	"coffee-tracker-backend/internal/infrastructure/utils"
	"coffee-tracker-backend/internal/repositories"

	"github.com/google/uuid"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) repositories.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Create inserts a new user
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, email, mobile, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		utils.SafeToLower(user.Email),
		utils.NullIfEmpty(user.Mobile),
		utils.NullIfEmpty(user.Name),
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// getUserByField is a helper for fetching users by any field
func (r *UserRepositoryImpl) getUserByField(ctx context.Context, field string, value interface{}) (*entities.User, error) {
	query := fmt.Sprintf(`
		SELECT id, email, mobile, name, COALESCE(avatar_url, '') AS avatar_url, status_id, created_at, updated_at
		FROM users
		WHERE %s = $1
	`, field)

	var user entities.User
	err := r.db.QueryRowContext(ctx, query, value).Scan(
		&user.ID, &user.Email, &user.Mobile, &user.Name,
		&user.AvatarURL, &user.StatusID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found by %s=%v: %w", field, value, err)
	}
	return &user, nil
}

// GetByID fetches a user by ID
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return r.getUserByField(ctx, "id", id)
}

// GetByEmail fetches a user by email
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	return r.getUserByField(ctx, "email", email)
}

// GetByMobile fetches a user by mobile
func (r *UserRepositoryImpl) GetByMobile(ctx context.Context, mobile string) (*entities.User, error) {
	return r.getUserByField(ctx, "mobile", mobile)
}

// Update updates user's email, mobile, name
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users 
		SET email = $2, mobile = $3, name = $4, updated_at = $5
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		utils.SafeToLower(user.Email),
		utils.NullIfEmpty(user.Mobile),
		utils.NullIfEmpty(user.Name),
		utils.NowUTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to update user %s: %w", user.ID, err)
	}
	return nil
}

// Delete removes a user by ID
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user %s: %w", id, err)
	}
	return nil
}

// UpdateProfile updates user profile fields based on request DTO
func (r *UserRepositoryImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, req *models.UpdateUserProfileRequest) error {
	query := `UPDATE users SET `
	params := []interface{}{}
	i := 1

	if req.Name != nil {
		query += `name = $` + strconv.Itoa(i) + `, `
		params = append(params, utils.NullIfEmpty(*req.Name))
		i++
	}
	if req.Email != nil {
		query += `email = $` + strconv.Itoa(i) + `, `
		params = append(params, utils.SafeToLower(*req.Email))
		i++
	}
	// if req.Address != nil {
	// 	query += `address = $` + strconv.Itoa(i) + `, `
	// 	params = append(params, *req.Address)
	// 	i++
	// }
	// if req.City != nil {
	// 	query += `city = $` + strconv.Itoa(i) + `, `
	// 	params = append(params, *req.City)
	// 	i++
	// }
	// if req.ZipCode != nil {
	// 	query += `zip_code = $` + strconv.Itoa(i) + `, `
	// 	params = append(params, *req.ZipCode)
	// 	i++
	// }

	// Remove trailing comma
	query = strings.TrimSuffix(query, ", ")
	// Always update updated_at
	query = strings.TrimSuffix(query, ", ")
	query += `, updated_at = $` + strconv.Itoa(i) + ` WHERE id = $` + strconv.Itoa(i+1)
	params = append(params, utils.NowUTC(), userID)

	_, err := r.db.ExecContext(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("failed to update profile for user %s: %w", userID, err)
	}
	return nil
}

// UpdateProfileImage updates user's avatar_url
func (r *UserRepositoryImpl) UpdateProfileImage(ctx context.Context, user *entities.User) error {
	query := `UPDATE users SET avatar_url = $1, updated_at = $3 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, user.AvatarURL, user.ID, utils.NowUTC())
	if err != nil {
		return fmt.Errorf("failed to update profile image for user %s: %w", user.ID, err)
	}
	return nil
}

// DeleteProfileImage sets user's avatar_url to NULL
func (r *UserRepositoryImpl) DeleteProfileImage(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET avatar_url = NULL, updated_at = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID, utils.NowUTC())
	if err != nil {
		return fmt.Errorf("failed to delete profile image for user %s: %w", userID, err)
	}
	return nil
}
