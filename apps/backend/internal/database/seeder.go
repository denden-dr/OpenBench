package database

import (
	"context"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SeedDevUsers seeds the initial admin user if it doesn't already exist
func SeedDevUsers(db *Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	queryCheck := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := db.DB.GetContext(ctx, &exists, queryCheck, "admin@openbench.dev")
	if err != nil {
		return err
	}

	if exists {
		log.Println("Admin seeder: user admin@openbench.dev already exists, skipping")
		return nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("SecureAdminPassword123!"), 12)
	if err != nil {
		return err
	}

	queryInsert := "INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3)"
	_, err = db.DB.ExecContext(ctx, queryInsert, "admin@openbench.dev", string(passwordHash), "admin")
	if err != nil {
		return err
	}

	log.Println("Admin seeder: user admin@openbench.dev seeded successfully")
	return nil
}
