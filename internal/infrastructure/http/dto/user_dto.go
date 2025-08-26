package dto

import (
	"github.com/google/uuid"
)

type LoggedInUserResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Mobile       string    `json:"mobile"`
}

// type UpdateUserRequest struct {
// 	Name  string  `json:"name" binding:"required,min=2,max=100"`
// 	Email *string `json:"email,omitempty" binding:"omitempty,email"`
// }