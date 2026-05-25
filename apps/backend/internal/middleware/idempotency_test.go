package middleware

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestIdempotencyConcretePath(t *testing.T) {
	app := fiber.New()

	var actualPath string
	var actualOk bool

	app.Use(func(c *fiber.Ctx) error {
		actualPath, actualOk = idempotencyConcretePath(c)
		return c.SendStatus(fiber.StatusOK)
	})

	tests := []struct {
		name       string
		method     string
		path       string
		expectPath string
		expectOk   bool
	}{
		{
			name:       "valid POST tickets",
			method:     fiber.MethodPost,
			path:       "/api/v1/tickets",
			expectPath: "/api/v1/tickets",
			expectOk:   true,
		},
		{
			name:       "valid POST tickets with trailing slash",
			method:     fiber.MethodPost,
			path:       "/api/v1/tickets/",
			expectPath: "/api/v1/tickets",
			expectOk:   true,
		},
		{
			name:       "valid PATCH ticket ID",
			method:     fiber.MethodPatch,
			path:       "/api/v1/tickets/123",
			expectPath: "/api/v1/tickets/123",
			expectOk:   true,
		},
		{
			name:       "valid PATCH ticket ID trailing slash",
			method:     fiber.MethodPatch,
			path:       "/api/v1/tickets/123/",
			expectPath: "/api/v1/tickets/123",
			expectOk:   true,
		},
		{
			name:       "invalid method GET tickets",
			method:     fiber.MethodGet,
			path:       "/api/v1/tickets",
			expectPath: "",
			expectOk:   false,
		},
		{
			name:       "invalid path nested patch",
			method:     fiber.MethodPatch,
			path:       "/api/v1/tickets/123/subroute",
			expectPath: "",
			expectOk:   false,
		},
		{
			name:       "empty ID on PATCH",
			method:     fiber.MethodPatch,
			path:       "/api/v1/tickets/",
			expectPath: "",
			expectOk:   false,
		},
		{
			name:       "valid POST warranty claims",
			method:     fiber.MethodPost,
			path:       "/api/v1/warranty-claims",
			expectPath: "/api/v1/warranty-claims",
			expectOk:   true,
		},
		{
			name:       "valid POST warranty claims with trailing slash",
			method:     fiber.MethodPost,
			path:       "/api/v1/warranty-claims/",
			expectPath: "/api/v1/warranty-claims",
			expectOk:   true,
		},
		{
			name:       "valid POST warranty claims approve",
			method:     fiber.MethodPost,
			path:       "/api/v1/warranty-claims/123/approve",
			expectPath: "/api/v1/warranty-claims/123/approve",
			expectOk:   true,
		},
		{
			name:       "valid POST warranty claims approve trailing slash",
			method:     fiber.MethodPost,
			path:       "/api/v1/warranty-claims/123/approve/",
			expectPath: "/api/v1/warranty-claims/123/approve",
			expectOk:   true,
		},
		{
			name:       "valid POST warranty claims void",
			method:     fiber.MethodPost,
			path:       "/api/v1/warranty-claims/456/void",
			expectPath: "/api/v1/warranty-claims/456/void",
			expectOk:   true,
		},
		{
			name:       "invalid method GET warranty claims",
			method:     fiber.MethodGet,
			path:       "/api/v1/warranty-claims",
			expectPath: "",
			expectOk:   false,
		},
		{
			name:       "invalid path nested warranty claim subroute",
			method:     fiber.MethodPost,
			path:       "/api/v1/warranty-claims/123/approve/extra",
			expectPath: "",
			expectOk:   false,
		},
		{
			name:       "invalid path wrong action word",
			method:     fiber.MethodPost,
			path:       "/api/v1/warranty-claims/123/something",
			expectPath: "",
			expectOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualPath = ""
			actualOk = false

			req, err := http.NewRequest(tt.method, tt.path, nil)
			assert.NoError(t, err)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			resp.Body.Close()

			assert.Equal(t, tt.expectOk, actualOk)
			assert.Equal(t, tt.expectPath, actualPath)
		})
	}
}

func TestValidateScopedIdempotencyKey(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		expectErr bool
	}{
		{
			name:      "valid scoped key",
			key:       "POST:/api/v1/tickets:6f37c5c6-2a74-4b8b-9b85-10efc9b9122d",
			expectErr: false,
		},
		{
			name:      "valid patch scoped key",
			key:       "PATCH:/api/v1/tickets/123:6f37c5c6-2a74-4b8b-9b85-10efc9b9122d",
			expectErr: false,
		},
		{
			name:      "too few parts",
			key:       "tickets:6f37c5c6-2a74-4b8b-9b85-10efc9b9122d",
			expectErr: true,
		},
		{
			name:      "invalid UUID suffix",
			key:       "POST:/api/v1/tickets:invalid-uuid",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateScopedIdempotencyKey(tt.key)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHashIdempotencyRequest(t *testing.T) {
	method := "POST"
	path := "/api/v1/tickets"
	body := []byte(`{"customer_name":"Budi"}`)

	hash1 := hashIdempotencyRequest(method, path, body)
	hash2 := hashIdempotencyRequest(method, path, body)

	assert.NotEmpty(t, hash1)
	assert.Equal(t, hash1, hash2, "hashes should be stable and identical for identical inputs")

	hashDiffBody := hashIdempotencyRequest(method, path, []byte(`{"customer_name":"Andi"}`))
	assert.NotEqual(t, hash1, hashDiffBody, "hashes should differ when request body changes")

	hashDiffMethod := hashIdempotencyRequest("PATCH", path, body)
	assert.NotEqual(t, hash1, hashDiffMethod, "hashes should differ when HTTP method changes")
}

func TestHashIdempotencyRequestChangesAfterValidationCorrection(t *testing.T) {
	method := "POST"
	path := "/api/v1/tickets"
	invalidBody := []byte(`{"customer_name":"Budi","customer_gender":"Male"}`)
	correctedBody := []byte(`{"customer_name":"Budi","customer_gender":"Male","brand":"Apple","model":"iPhone 13","issue":"LCD Mati"}`)

	invalidHash := hashIdempotencyRequest(method, path, invalidBody)
	correctedHash := hashIdempotencyRequest(method, path, correctedBody)

	assert.NotEqual(t, invalidHash, correctedHash)
}
