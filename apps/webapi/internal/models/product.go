package models

import "time"

type Product struct {
	ID        string     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Price     int64      `json:"price" db:"price"`
	Stock     int        `json:"stock" db:"stock"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}
