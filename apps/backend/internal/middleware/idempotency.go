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

func idempotencyConcretePath(c *fiber.Ctx) (string, bool) {
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

	if method == fiber.MethodPost && path == "/api/v1/warranty-claims" {
		return "/api/v1/warranty-claims", true
	}

	if method == fiber.MethodPost && strings.HasPrefix(path, "/api/v1/warranty-claims/") {
		suffix := strings.TrimPrefix(path, "/api/v1/warranty-claims/")
		parts := strings.Split(suffix, "/")
		if len(parts) == 2 && (parts[1] == "approve" || parts[1] == "void") {
			id := parts[0]
			if id != "" {
				return "/api/v1/warranty-claims/" + id + "/" + parts[1], true
			}
		}
	}

	return "", false
}

func ScopeIdempotencyKey(store *database.PostgresStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Clean spoofed internal header from client
		c.Request().Header.Del(scopedIdempotencyHeader)

		concretePath, ok := idempotencyConcretePath(c)
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

func NewIdempotency(storage fiber.Storage) fiber.Handler {
	return idempotency.New(idempotency.Config{
		Next: func(c *fiber.Ctx) bool {
			_, ok := idempotencyConcretePath(c)
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
