package auth

import "time"

type User struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type RefreshTokenRecord struct {
	ID                string     `db:"id"`
	FamilyID          string     `db:"family_id"`
	UserID            string     `db:"user_id"`
	TokenHash         string     `db:"token_hash"`
	IsUsed            bool       `db:"is_used"`
	IsRevoked         bool       `db:"is_revoked"`
	UsedAt            *time.Time `db:"used_at"`
	ExpiresAt         time.Time  `db:"expires_at"`
	CreatedAt         time.Time  `db:"created_at"`
	ReplacedByTokenID *string    `db:"replaced_by_token_id"`
}
