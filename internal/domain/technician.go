package domain

import (
	"time"

	"github.com/google/uuid"
)

// Technician represents the profile for a user with tech capabilities.
type Technician struct {
	UserID                uuid.UUID `db:"user_id" json:"user_id"`
	Bio                   *string   `db:"bio" json:"bio,omitempty"`
	Specialties           *string   `db:"specialties" json:"specialties,omitempty"`
	Rating                float64   `db:"rating" json:"rating"`
	TotalRepairsCompleted int       `db:"total_repairs_completed" json:"total_repairs_completed"`
	IsActive              bool      `db:"is_active" json:"is_active"`
	CreatedAt             time.Time `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time `db:"updated_at" json:"updated_at"`
}
