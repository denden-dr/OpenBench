package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

const (
	DefaultLimit = 10
	MaxLimit     = 25
)

type CursorMeta struct {
	NextCursor *string `json:"next_cursor"`
	Limit      int     `json:"limit"`
}

type CursorPaginatedResponse[T any] struct {
	Data []T        `json:"data"`
	Meta CursorMeta `json:"meta"`
}

// NewCursorPaginatedResponse creates a standardized cursor paginated response.
func NewCursorPaginatedResponse[T any](data []T, limit int, nextCursor string) CursorPaginatedResponse[T] {
	if data == nil {
		data = make([]T, 0)
	}

	var next *string
	if nextCursor != "" {
		next = &nextCursor
	}

	return CursorPaginatedResponse[T]{
		Data: data,
		Meta: CursorMeta{
			NextCursor: next,
			Limit:      limit,
		},
	}
}

// ParseCursorPagination parses 'limit' and 'cursor' query parameters from Fiber context.
// 'limit' is capped at MaxLimit (25).
func ParseCursorPagination(c fiber.Ctx) (limit int, cursor string) {
	limitStr := c.Query("limit")
	limit = DefaultLimit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}

	cursor = c.Query("cursor")
	return limit, cursor
}
