package repository

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrClaimConflict = errors.New("ticket is already claimed")
)
