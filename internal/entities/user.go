// file: internal/entities/user.go
package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type User struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Email     string     `db:"email" json:"email"`
	Mobile    string     `db:"mobile" json:"mobile"`
	Name      string     `db:"name" json:"name"`
	AvatarURL string     `db:"avatar_url" json:"avatar_url"`
	StatusID  int        `db:"status_id" json:"status_id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}
const (
	StatusPending   = 1
	StatusActive    = 2
	StatusSuspended = 3
	StatusInactive  = 4
	StatusBanned    = 5
	StatusDeleted   = 6
	StatusArchived  = 7
)

// Business logic check
func (u User) IsActive() bool {
	return u.StatusID == StatusActive
}