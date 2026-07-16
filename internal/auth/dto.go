package auth

import "time"

type SuccessResponse[T any] struct {
	Data T `json:"data"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type UserProfileResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginResponse struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"-"` // Not exposed in JSON, used for cookie
	ExpiresAt    time.Time           `json:"expires_at"`
	User         UserProfileResponse `json:"user"`
}

type RefreshResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"-"`
	ExpiresAt    time.Time `json:"expires_at"`
}
