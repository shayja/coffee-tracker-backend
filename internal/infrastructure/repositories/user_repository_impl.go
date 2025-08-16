package repositories

import (
	"context"
	"database/sql"

	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"

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
		INSERT INTO users (id, email, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.CreatedAt,
		user.UpdatedAt,
	)
	
	return err
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	query := `
		SELECT id, email, name, status_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	
	var user entities.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
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
		SELECT id, email, name, status_id, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	
	var user entities.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
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
		SET email = $2, name = $3, updated_at = $4
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
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

func (r *AuthRepositoryImpl) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error) {
    var userID uuid.UUID
    query := `
        SELECT user_id
        FROM user_refresh_tokens
        WHERE token = $1 AND expires_at > NOW()
        LIMIT 1
    `
    err := r.db.QueryRowContext(ctx, query, refreshToken).Scan(&userID)
    if err != nil {
        if err == sql.ErrNoRows {
            return uuid.Nil, repositories.ErrNotFound
        }
        return uuid.Nil, err
    }
    return userID, nil
}