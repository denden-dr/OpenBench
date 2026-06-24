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

// SeedDevProducts seeds initial products in development/test environments
func SeedDevProducts(db *Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	queryCheck := "SELECT COUNT(*) FROM products"
	err := db.DB.GetContext(ctx, &count, queryCheck)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Product seeder: products already exist, skipping")
		return nil
	}

	products := []struct {
		ID        string
		Name      string
		Category  string
		Stock     int
		Price     float64
		CostPrice float64
		MinStock  int
	}{
		{"550e8400-e29b-41d4-a716-446655440011", "Charger 25W Fast Charging Type-C", "retail", 18, 245000, 150000, 5},
		{"550e8400-e29b-41d4-a716-446655440012", "Tempered Glass Ultra Clear iPhone 14 Pro", "retail", 3, 95000, 40000, 5},
		{"550e8400-e29b-41d4-a716-446655440013", "LCD Screen Module Samsung S23 Ultra (Original)", "spare_part", 4, 2600000, 1900000, 2},
		{"550e8400-e29b-41d4-a716-446655440014", "Battery Replacement iPhone 14 Pro (Original)", "spare_part", 12, 850000, 500000, 3},
		{"550e8400-e29b-41d4-a716-446655440015", "Soft Case Transparent anti-crack (Universal)", "retail", 25, 49000, 15000, 5},
	}

	queryInsert := `
		INSERT INTO products (id, name, category, stock, price, cost_price, min_stock)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	for _, p := range products {
		_, err := db.DB.ExecContext(ctx, queryInsert, p.ID, p.Name, p.Category, p.Stock, p.Price, p.CostPrice, p.MinStock)
		if err != nil {
			return err
		}
	}

	log.Println("Product seeder: 5 products seeded successfully")
	return nil
}

