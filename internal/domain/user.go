package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents the local proxy of a Supabase authenticated user.
type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	FullName  *string   `db:"full_name" json:"full_name,omitempty"`
	AvatarURL *string   `db:"avatar_url" json:"avatar_url,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
