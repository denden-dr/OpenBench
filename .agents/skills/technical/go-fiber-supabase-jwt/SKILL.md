---
name: go-fiber-supabase-jwt
description: Use when implementing authentication middleware, protecting API routes, extracting user identity from JWTs in cookies, or setting up role-based access control in the Go Fiber backend with Supabase Auth.
---

# Go Fiber + Supabase Cookie-Based Authentication

## Overview

Supabase JWTs are stored in `HttpOnly` cookies. The Go Fiber backend extracts the token from the `Cookie` header, verifies it, and identifies the user.

## Auth Flow

```
Client → Go Fiber API
  1. Extract sb-access-token from Cookie header
  2. Verify JWT signature using Supabase JWT secret (HMAC, local)
  3. Extract supabase_uid from "sub" claim
  4. Look up internal user by supabase_uid → get user_id + role
  5. Set user_id, supabase_uid, role in c.Locals()
```

## Middleware Pattern (Cookie Extraction)

```go
func (m *AuthMiddleware) Authenticate() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Extract from cookie
        tokenString := c.Cookies("sb-access-token")
        if tokenString == "" {
            return c.Status(401).JSON(ErrorResponse{Code: "UNAUTHORIZED"})
        }

        token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method")
            }
            return m.jwtSecret, nil
        })

        if err != nil || !token.Valid {
            return c.Status(401).JSON(ErrorResponse{Code: "UNAUTHORIZED"})
        }

        claims := token.Claims.(jwt.MapClaims)
        supabaseUID, _ := claims["sub"].(string)

        userID, role, err := m.userLookup(c.Context(), supabaseUID)
        if err != nil {
            return c.Status(401).JSON(ErrorResponse{Code: "UNAUTHORIZED", Message: "User not found"})
        }

        c.Locals("user_id", userID)
        c.Locals("supabase_uid", supabaseUID)
        c.Locals("role", role)
        return c.Next()
    }
}
```

## Security: CSRF Protection

When using cookies, verify the `Origin` or `Referer` to prevent CSRF.

```go
func CSRFMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        if c.Method() == "GET" || c.Method() == "HEAD" || c.Method() == "OPTIONS" {
            return c.Next()
        }
        
        origin := c.Get("Origin")
        if origin != os.Getenv("FRONTEND_URL") {
            return c.Status(403).JSON(ErrorResponse{Code: "FORBIDDEN_ORIGIN"})
        }
        return c.Next()
    }
}
```

## Role-Based Authorization

```go
func RequireRole(roles ...string) fiber.Handler {
    allowed := make(map[string]bool, len(roles))
    for _, r := range roles { allowed[r] = true }
    return func(c *fiber.Ctx) error {
        role, ok := c.Locals("role").(string)
        if !ok || !allowed[role] {
            return c.Status(403).JSON(ErrorResponse{Code: "FORBIDDEN"})
        }
        return c.Next()
    }
}
```

## Context Locals Reference

| Key | Type | Contains |
|-----|------|----------|
| `user_id` | string | Internal UUID from `users` table |
| `supabase_uid` | string | Supabase `auth.users.id` |
| `role` | string | `admin`, `technician`, or `customer` |

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Expecting `Authorization: Bearer` header | Use `c.Cookies("sb-access-token")` |
| Forgetting CSRF protection | Cookies are sent automatically; check `Origin` header |
| Not setting `SameSite=Lax` in frontend | Required for cookies to be sent across subdomains |
| Hardcoding JWT secret | Load from env |
