package auth

import "time"

// SignInRequest represents the expected payload for the sign-in endpoint
type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// SignInResult represents the result of a successful service-level sign-in
type SignInResult struct {
	User            *User
	AccessToken     string
	RawRefreshToken string
	ExpiresAt       time.Time
}

// RefreshResult represents the result of a successful service-level refresh
type RefreshResult struct {
	AccessToken     string
	RawRefreshToken string
	ExpiresAt       time.Time
	GraceRefreshed  bool
}

// SignInResponse is the response payload sent back inside response.Data
type SignInResponse struct {
	Role   string `json:"role"`
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// MeResponse is the response payload sent back inside response.Data for the /me endpoint
type MeResponse struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
}
