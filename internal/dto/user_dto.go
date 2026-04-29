package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse defines the public-facing user profile information.
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  *string   `json:"full_name,omitempty"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}
