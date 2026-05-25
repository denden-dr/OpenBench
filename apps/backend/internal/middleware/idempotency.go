package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/google/uuid"
)

const (
	rawIdempotencyHeader    = "X-Idempotency-Key"
	scopedIdempotencyHeader = "X-Scoped-Idempotency-Key"
)

func ticketIdempotencyConcretePath(c *fiber.Ctx) (string, bool) {
	path := strings.TrimRight(c.Path(), "/")
	method := c.Method()

	if method == fiber.MethodPost && path == "/api/v1/tickets" {
		return "/api/v1/tickets", true
	}

	if method == fiber.MethodPatch && strings.HasPrefix(path, "/api/v1/tickets/") {
		id := strings.TrimPrefix(path, "/api/v1/tickets/")
		if id != "" && !strings.Contains(id, "/") {
			return "/api/v1/tickets/" + id, true
		}
	}

	return "", false
}

func ScopeTicketIdempotencyKey(store *database.PostgresStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Clean spoofed internal header from client
		c.Request().Header.Del(scopedIdempotencyHeader)

		concretePath, ok := ticketIdempotencyConcretePath(c)
		if !ok {
			return c.Next()
		}

		rawKey := c.Get(rawIdempotencyHeader)
		if rawKey == "" {
			return c.Next()
		}

		if _, err := uuid.Parse(rawKey); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid idempotency key")
		}

		scopedKey := fmt.Sprintf("%s:%s:%s", c.Method(), concretePath, rawKey)
		requestHash := hashIdempotencyRequest(c.Method(), concretePath, c.BodyRaw())
		if err := store.ReserveRequest(scopedKey, requestHash, 30*time.Minute); err != nil {
			if errors.Is(err, database.ErrIdempotencyConflict) {
				return fiber.NewError(fiber.StatusConflict, "idempotency key reused with different request body")
			}
			return err
		}

		c.Request().Header.Set(scopedIdempotencyHeader, scopedKey)
		return c.Next()
	}
}

func hashIdempotencyRequest(method string, concretePath string, body []byte) string {
	h := sha256.New()
	h.Write([]byte(method))
	h.Write([]byte("\n"))
	h.Write([]byte(concretePath))
	h.Write([]byte("\n"))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

func NewTicketIdempotency(storage fiber.Storage) fiber.Handler {
	return idempotency.New(idempotency.Config{
		Next: func(c *fiber.Ctx) bool {
			_, ok := ticketIdempotencyConcretePath(c)
			return !ok
		},
		Storage:           storage,
		Lifetime:          30 * time.Minute,
		KeyHeader:         scopedIdempotencyHeader,
		KeyHeaderValidate: validateScopedIdempotencyKey,
	})
}

func validateScopedIdempotencyKey(key string) error {
	parts := strings.Split(key, ":")
	if len(parts) < 3 {
		return idempotency.ErrInvalidIdempotencyKey
	}

	rawUUID := parts[len(parts)-1]
	if _, err := uuid.Parse(rawUUID); err != nil {
		return fmt.Errorf("%w: %v", idempotency.ErrInvalidIdempotencyKey, err)
	}

	return nil
}
