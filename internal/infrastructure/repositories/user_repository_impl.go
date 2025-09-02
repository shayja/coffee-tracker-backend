// file: internal/infrastructure/repositories/user_repository_impl.go
package repositories

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"
	"coffee-tracker-backend/internal/infrastructure/http/dto"

	"github.com/google/uuid"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) repositories.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, email, mobile, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Mobile,
		user.Name,
		user.CreatedAt,
		user.UpdatedAt,
	)
	
	return err
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	query := `
		SELECT id, email, mobile, name, COALESCE(avatar_url, '') AS avatar_url, status_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user entities.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Mobile,
		&user.Name,
		&user.AvatarURL,
		&user.StatusID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, email, mobile, name, COALESCE(avatar_url, '') AS avatar_url, status_id, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user entities.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Mobile,
		&user.Name,
		&user.AvatarURL,
		&user.StatusID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByMobile(ctx context.Context, mobile string) (*entities.User, error) {
	query := `
		SELECT id, email, mobile, name, COALESCE(avatar_url, '') AS avatar_url, status_id, created_at, updated_at
		FROM users
		WHERE mobile = $1
	`
	var user entities.User
	err := r.db.QueryRowContext(ctx, query, mobile).Scan(
		&user.ID,
		&user.Email,
		&user.Mobile,
		&user.Name,
		&user.AvatarURL,
		&user.StatusID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}


func (r *UserRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users 
		SET email = $2, mobile= $3, name = $4, updated_at = $5
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Mobile,
		user.Name,
		user.UpdatedAt,
	)
	
	return err
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// UpdateProfile updates the user's profile fields (name, avatar_url).
func (r *UserRepositoryImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserProfileRequest) error {
	query := `UPDATE users SET `
	params := []interface{}{}
	i := 1

	if req.Name != nil {
		query += `name = $` + strconv.Itoa(i) + `, `
		params = append(params, *req.Name)
		i++
	}
	if req.Email != nil {
		query += `email = $` + strconv.Itoa(i) + `, `
		params = append(params, *req.Email)
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

	// Always update updated_at
	query += `updated_at = NOW() WHERE id = $` + strconv.Itoa(i)
	params = append(params, userID)

	_, err := r.db.ExecContext(ctx, query, params...)
	return err
}

func (r *UserRepositoryImpl) UpdateAProfileImage(ctx context.Context, user *entities.User) error {
    query := `UPDATE users SET avatar_url = $1, updated_at = $2 WHERE id = $3`
    _, err := r.db.ExecContext(ctx, query, user.AvatarURL, user.UpdatedAt, user.ID)
    return err
}


func (r *UserRepositoryImpl) DeleteProfileImage(ctx context.Context, userID uuid.UUID, updatedAt time.Time) error {
	query := `
		UPDATE users
		SET avatar_url = NULL, updated_at = $2
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, userID, updatedAt)
	return err
}