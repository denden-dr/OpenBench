---
name: openbench-api-patterns
description: Use when creating API endpoints, writing handlers, defining request/response DTOs, implementing error handling, or structuring the Go Fiber backend service layer in the OpenBench project
---

# OpenBench API Patterns

## Overview

Standard patterns for the Go Fiber REST API backend. Every endpoint must follow these conventions for consistency. **Consult `openbench-domain-guide` for business rules before implementing any endpoint.**

## Endpoint Naming

```
GET    /api/v1/{resource}          → List (paginated)
GET    /api/v1/{resource}/:id      → Get one
POST   /api/v1/{resource}          → Create
PUT    /api/v1/{resource}/:id      → Full update
PATCH  /api/v1/{resource}/:id      → Partial update
DELETE /api/v1/{resource}/:id      → Delete (soft)

# Actions (non-CRUD)
POST   /api/v1/tickets/:id/claim   → Technician claims ticket
POST   /api/v1/tickets/:id/parts   → Log parts used
POST   /api/v1/tickets/:id/status  → Update status

# Webhooks (external)
POST   /api/v1/webhooks/payment    → Midtrans callback
```

**Rules:**
- Always plural nouns for resources (`/tickets` not `/ticket`)
- Nest only one level deep (`/tickets/:id/parts`, never `/tickets/:id/parts/:partId/logs`)
- Action endpoints use `POST` with verb in path
- API version prefix: `/api/v1/`

## Response Envelope

All responses use a consistent envelope:

```go
// Success response
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

// Error response
type ErrorResponse struct {
    Success bool         `json:"success"`
    Error   ErrorDetail  `json:"error"`
}

type ErrorDetail struct {
    Code    string            `json:"code"`
    Message string            `json:"message"`
    Details map[string]string `json:"details,omitempty"` // field-level validation
}

// Pagination meta
type Meta struct {
    Page       int `json:"page"`
    PerPage    int `json:"per_page"`
    Total      int64 `json:"total"`
    TotalPages int `json:"total_pages"`
}
```

### Example Responses

```json
// Success — single resource
{ "success": true, "data": { "id": "abc-123", "status": "received" } }

// Success — list with pagination
{ "success": true, "data": [...], "meta": { "page": 1, "per_page": 20, "total": 45, "total_pages": 3 } }

// Validation error
{ "success": false, "error": { "code": "VALIDATION_ERROR", "message": "Invalid input", "details": { "phone": "must be a valid phone number", "device_type": "must be 'Android' or 'Apple'" } } }

// Business logic error
{ "success": false, "error": { "code": "TICKET_ALREADY_CLAIMED", "message": "This ticket has already been assigned to another technician" } }

// Auth error
{ "success": false, "error": { "code": "UNAUTHORIZED", "message": "Invalid or expired token" } }
```

## Error Codes

| Code | HTTP | When |
|------|------|------|
| `VALIDATION_ERROR` | 400 | Request body fails validation |
| `INSUFFICIENT_STOCK` | 400 | Part stock too low for requested quantity |
| `INVALID_STATUS_TRANSITION` | 400 | Ticket status change not allowed |
| `UNAUTHORIZED` | 401 | Missing or invalid JWT |
| `FORBIDDEN` | 403 | User role lacks permission |
| `NOT_FOUND` | 404 | Resource doesn't exist |
| `TICKET_ALREADY_CLAIMED` | 409 | Technician tried to claim an already-assigned ticket |
| `DUPLICATE_ENTRY` | 409 | Unique constraint violation |
| `PAYMENT_VERIFICATION_FAILED` | 422 | Webhook signature invalid |
| `INTERNAL_ERROR` | 500 | Unexpected server error (log details, return generic message) |

## Handler Structure

Every handler follows this pattern:

```go
func (h *TicketHandler) Create(c *fiber.Ctx) error {
    // 1. Extract auth context (set by middleware)
    userID := c.Locals("user_id").(string)
    role := c.Locals("role").(string)

    // 2. Parse & validate request
    var req CreateTicketRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(ErrorResponse{
            Error: ErrorDetail{Code: "VALIDATION_ERROR", Message: "Invalid request body"},
        })
    }
    if errs := validate(req); errs != nil {
        return c.Status(400).JSON(ErrorResponse{
            Error: ErrorDetail{Code: "VALIDATION_ERROR", Message: "Invalid input", Details: errs},
        })
    }

    // 3. Call service (business logic lives here, NOT in handler)
    ticket, err := h.service.CreateTicket(c.Context(), userID, req)
    if err != nil {
        return mapServiceError(c, err) // maps domain errors to HTTP
    }

    // 4. Return response
    return c.Status(201).JSON(Response{Success: true, Data: ticket})
}
```

**Rules:**
- Handlers are **thin** — parse request, call service, return response
- **All business logic** lives in the service layer
- Services return domain errors, handlers map them to HTTP
- Always use `c.Context()` to propagate context for cancellation/timeouts

## Service Layer

```go
// Service interface
type TicketService interface {
    CreateTicket(ctx context.Context, userID string, req CreateTicketRequest) (*TicketDTO, error)
    ClaimTicket(ctx context.Context, techID string, ticketID string) (*TicketDTO, error)
    UpdateStatus(ctx context.Context, userID string, ticketID string, status string) (*TicketDTO, error)
}

// Implementation
type ticketService struct {
    repo     TicketRepository
    partRepo PartRepository
    audit    AuditService
}
```

**Rules:**
- Services accept and return **DTOs**, not database models
- Services own **transaction boundaries** (not repositories)
- Services call the audit service for sensitive operations
- Repository methods are simple CRUD — no business logic

## Pagination

Query parameters: `?page=1&per_page=20`

```go
// Defaults
const (
    DefaultPage    = 1
    DefaultPerPage = 20
    MaxPerPage     = 100
)
```

## Request Validation

Use struct tags with a validator (e.g., `go-playground/validator`):

```go
type CreateTicketRequest struct {
    DeviceType      string   `json:"device_type" validate:"required,oneof=Android Apple"`
    Brand           string   `json:"brand" validate:"required,max=100"`
    Model           string   `json:"model" validate:"required,max=100"`
    IssueDescription string  `json:"issue_description" validate:"required,max=2000"`
    DeviceAccessInfo string  `json:"device_access_info" validate:"required,max=500"`
    Accessories     []string `json:"accessories" validate:"dive,max=100"`
    TermsAgreed     bool     `json:"terms_agreed" validate:"required,eq=true"`
}
```

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Business logic in handler | Move to service layer |
| Raw SQL errors in response | Map to domain errors, log the SQL error |
| Missing audit log on sensitive action | Always call `audit.Log()` for status changes, stock updates, payment changes |
| Forgetting `c.Context()` | Always propagate context for cancellation |
| Inconsistent error code strings | Use constants from the error codes table above |
