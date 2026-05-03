# Full Business Cycle Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the complete repair lifecycle (Booking → Diagnosis → Approval/Cancellation → Repair → Ready → Pickup) with simulated User/Technician roles and no formal authentication.

**Architecture:**
- **Role Simulation**: A `RoleMiddleware` extracts roles from the `X-Mock-Role` header (`user` | `technician`). Default is `user`.
- **Domain Actions**: The service layer exposes domain-level methods (`ClaimTicket`, `CompleteDiagnosis`, `ApproveRepair`, etc.) instead of a generic `UpdateStatus(role, newStatus)`. Each method encodes the valid transition internally. The handler maps the HTTP role to the correct method call.
- **Claiming**: `ClaimTicket` is a dedicated operation that atomically sets `status = diagnosing` AND `technician_id` in one SQL statement. It is **not** handled by a generic status update.

**Tech Stack:** Go, Fiber, PostgreSQL, sqlx.

**Status Lifecycle:**

```
received → diagnosing → pending_approval → repairing → ready → picked_up
                                         ↘ cancelled → picked_up
```

| From | To | Who | Method |
|---|---|---|---|
| `received` | `diagnosing` | Technician | `ClaimTicket` |
| `diagnosing` | `pending_approval` | Technician | `CompleteDiagnosis` |
| `pending_approval` | `repairing` | User | `ApproveRepair` |
| `pending_approval` | `cancelled` | User | `CancelRepair` |
| `repairing` | `ready` | Technician | `CompleteRepair` |
| `ready` | `picked_up` | Technician | `MarkPickedUp` |
| `cancelled` | `picked_up` | Technician | `MarkPickedUp` |

---

### Task 1: Migration, Model, and DTOs

**Files:**
- Create: `apps/backend/migrations/000003_add_technician_to_tickets.up.sql`
- Create: `apps/backend/migrations/000003_add_technician_to_tickets.down.sql`
- Create: `apps/backend/internal/model/status.go`
- Modify: `apps/backend/internal/model/ticket.go`
- Modify: `apps/backend/internal/dto/ticket_dto.go`

- [ ] **Step 1: Create up migration**

```sql
-- apps/backend/migrations/000003_add_technician_to_tickets.up.sql
ALTER TABLE tickets ADD COLUMN technician_id UUID;
```

- [ ] **Step 2: Create down migration**

```sql
-- apps/backend/migrations/000003_add_technician_to_tickets.down.sql
ALTER TABLE tickets DROP COLUMN technician_id;
```

- [ ] **Step 3: Run migration**

Run: `make migrate-up`
Expected: Migration `000003` applied successfully.

- [ ] **Step 4: Create status constants**

```go
// apps/backend/internal/model/status.go
package model

const (
	StatusReceived        = "received"
	StatusDiagnosing      = "diagnosing"
	StatusPendingApproval = "pending_approval"
	StatusRepairing       = "repairing"
	StatusReady           = "ready"
	StatusCancelled       = "cancelled"
	StatusPickedUp        = "picked_up"
)

// ValidTransitions maps each status to its allowed next statuses.
var ValidTransitions = map[string][]string{
	StatusReceived:        {StatusDiagnosing},
	StatusDiagnosing:      {StatusPendingApproval},
	StatusPendingApproval: {StatusRepairing, StatusCancelled},
	StatusRepairing:       {StatusReady},
	StatusReady:           {StatusPickedUp},
	StatusCancelled:       {StatusPickedUp},
}
```

- [ ] **Step 5: Update Ticket model to include TechnicianID**

```go
// apps/backend/internal/model/ticket.go
// Add this field to the Ticket struct:
TechnicianID *string `db:"technician_id"`
```

- [ ] **Step 6: Add new DTOs**

```go
// apps/backend/internal/dto/ticket_dto.go
// Append these types to the existing file:

type UpdateTicketStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=repairing cancelled"`
}

type TicketBoardDTO struct {
	ID         string `json:"id"`
	DeviceType string `json:"device_type"`
	Brand      string `json:"brand"`
	Model      string `json:"model"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}
```

- [ ] **Step 7: Commit**

```bash
git add apps/backend/migrations/000003_*.sql apps/backend/internal/model/ apps/backend/internal/dto/
git commit -m "feat: add technician_id migration, status constants, and new DTOs"
```

---

### Task 2: Role Middleware

**Files:**
- Create: `apps/backend/internal/handler/middleware/role_middleware.go`

- [ ] **Step 1: Create role middleware**

```go
// apps/backend/internal/handler/middleware/role_middleware.go
package middleware

import "github.com/gofiber/fiber/v2"

const (
	RoleUser       = "user"
	RoleTechnician = "technician"
	LocalsRoleKey  = "role"
)

// RoleMiddleware extracts the X-Mock-Role header and stores it in Fiber locals.
// Defaults to "user" when the header is absent.
func RoleMiddleware(c *fiber.Ctx) error {
	role := c.Get("X-Mock-Role")
	if role != RoleTechnician {
		role = RoleUser
	}
	c.Locals(LocalsRoleKey, role)
	return c.Next()
}

// RequireTechnician rejects requests that don't have X-Mock-Role: technician.
func RequireTechnician(c *fiber.Ctx) error {
	if c.Locals(LocalsRoleKey) != RoleTechnician {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "Access denied: technician role required",
		})
	}
	return c.Next()
}
```

- [ ] **Step 2: Commit**

```bash
git add apps/backend/internal/handler/middleware/
git commit -m "feat: add mock role middleware"
```

---

### Task 3: Repository Layer

**Files:**
- Modify: `apps/backend/internal/repository/ticket_repo.go`
- Modify: `apps/backend/internal/repository/errors.go`

- [ ] **Step 1: Add new error sentinel**

```go
// apps/backend/internal/repository/errors.go
// Add to the existing var block:
ErrClaimConflict = errors.New("ticket is already claimed")
```

- [ ] **Step 2: Update TicketRepository interface**

Add these methods to the existing `TicketRepository` interface in `ticket_repo.go`:

```go
type TicketRepository interface {
	Create(ctx context.Context, ticket *model.Ticket) error
	GetByID(ctx context.Context, id string) (*model.Ticket, error)
	UpdateStatus(ctx context.Context, id string, newStatus string) error
	ClaimTicket(ctx context.Context, id string, technicianID string) error
	ListForBoard(ctx context.Context) ([]model.Ticket, error)
}
```

- [ ] **Step 3: Implement UpdateStatus**

```go
func (r *sqlTicketRepository) UpdateStatus(ctx context.Context, id string, newStatus string) error {
	query := `UPDATE tickets SET status = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, newStatus, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
```

- [ ] **Step 4: Implement ClaimTicket**

This atomically sets both `technician_id` and `status` only if the ticket is currently `received`. The `WHERE` clause prevents race conditions.

```go
func (r *sqlTicketRepository) ClaimTicket(ctx context.Context, id string, technicianID string) error {
	query := `
		UPDATE tickets
		SET status = $1, technician_id = $2
		WHERE id = $3 AND status = $4
	`
	result, err := r.db.ExecContext(ctx, query,
		model.StatusDiagnosing, technicianID, id, model.StatusReceived,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrClaimConflict
	}
	return nil
}
```

- [ ] **Step 5: Implement ListForBoard**

Returns only public-safe columns. Excludes `issue_description`, `diagnosis_fee`, and `technician_id`.

```go
func (r *sqlTicketRepository) ListForBoard(ctx context.Context) ([]model.Ticket, error) {
	var tickets []model.Ticket
	query := `
		SELECT id, device_type, brand, model, status, created_at
		FROM tickets
		ORDER BY created_at DESC
	`
	if err := r.db.SelectContext(ctx, &tickets, query); err != nil {
		return nil, err
	}
	return tickets, nil
}
```

- [ ] **Step 6: Commit**

```bash
git add apps/backend/internal/repository/
git commit -m "feat: add UpdateStatus, ClaimTicket, ListForBoard to repository"
```

---

### Task 4: Service Layer

**Files:**
- Modify: `apps/backend/internal/service/ticket_service.go`

- [ ] **Step 1: Add new error sentinels**

```go
// Add to the existing var block in ticket_service.go:
var (
	ErrTicketNotFound    = errors.New("ticket not found")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrClaimConflict     = errors.New("ticket is already claimed or not in received status")
)
```

- [ ] **Step 2: Update TicketService interface**

```go
type TicketService interface {
	CreateTicket(ctx context.Context, req *dto.CreateTicketRequest) (*dto.TicketResponse, error)
	GetTicket(ctx context.Context, id string) (*dto.TicketResponse, error)
	ClaimTicket(ctx context.Context, ticketID string, technicianID string) error
	CompleteDiagnosis(ctx context.Context, ticketID string) error
	ApproveRepair(ctx context.Context, ticketID string) error
	CancelRepair(ctx context.Context, ticketID string) error
	CompleteRepair(ctx context.Context, ticketID string) error
	MarkPickedUp(ctx context.Context, ticketID string) error
	ListForBoard(ctx context.Context) ([]dto.TicketBoardDTO, error)
}
```

- [ ] **Step 3: Implement ClaimTicket**

```go
func (s *ticketService) ClaimTicket(ctx context.Context, ticketID string, technicianID string) error {
	if err := s.repo.ClaimTicket(ctx, ticketID, technicianID); err != nil {
		if errors.Is(err, repository.ErrClaimConflict) {
			return ErrClaimConflict
		}
		return err
	}
	return nil
}
```

- [ ] **Step 4: Implement transition helpers**

Each method fetches the ticket, validates the current status, and calls `repo.UpdateStatus`.

```go
// transitionStatus is a private helper that validates and executes a status transition.
func (s *ticketService) transitionStatus(ctx context.Context, ticketID string, expectedCurrent string, newStatus string) error {
	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTicketNotFound
		}
		return err
	}
	if ticket.Status != expectedCurrent {
		return fmt.Errorf("%w: cannot move from '%s' to '%s'", ErrInvalidTransition, ticket.Status, newStatus)
	}
	return s.repo.UpdateStatus(ctx, ticketID, newStatus)
}

func (s *ticketService) CompleteDiagnosis(ctx context.Context, ticketID string) error {
	return s.transitionStatus(ctx, ticketID, model.StatusDiagnosing, model.StatusPendingApproval)
}

func (s *ticketService) ApproveRepair(ctx context.Context, ticketID string) error {
	return s.transitionStatus(ctx, ticketID, model.StatusPendingApproval, model.StatusRepairing)
}

func (s *ticketService) CancelRepair(ctx context.Context, ticketID string) error {
	return s.transitionStatus(ctx, ticketID, model.StatusPendingApproval, model.StatusCancelled)
}

func (s *ticketService) CompleteRepair(ctx context.Context, ticketID string) error {
	return s.transitionStatus(ctx, ticketID, model.StatusRepairing, model.StatusReady)
}

func (s *ticketService) MarkPickedUp(ctx context.Context, ticketID string) error {
	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTicketNotFound
		}
		return err
	}
	if ticket.Status != model.StatusReady && ticket.Status != model.StatusCancelled {
		return fmt.Errorf("%w: cannot pick up from '%s'", ErrInvalidTransition, ticket.Status)
	}
	return s.repo.UpdateStatus(ctx, ticketID, model.StatusPickedUp)
}
```

- [ ] **Step 5: Implement ListForBoard**

```go
func (s *ticketService) ListForBoard(ctx context.Context) ([]dto.TicketBoardDTO, error) {
	tickets, err := s.repo.ListForBoard(ctx)
	if err != nil {
		return nil, err
	}

	board := make([]dto.TicketBoardDTO, len(tickets))
	for i, t := range tickets {
		board[i] = dto.TicketBoardDTO{
			ID:         t.ID,
			DeviceType: t.DeviceType,
			Brand:      t.Brand,
			Model:      t.Model,
			Status:     t.Status,
			CreatedAt:  t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return board, nil
}
```

- [ ] **Step 6: Add `"fmt"` to the imports if not already present**

- [ ] **Step 7: Commit**

```bash
git add apps/backend/internal/service/
git commit -m "feat: add domain-level ticket lifecycle methods to service"
```

---

### Task 5: Handler Layer

**Files:**
- Modify: `apps/backend/internal/handler/ticket_handler.go`

- [ ] **Step 1: Add ClaimTicket handler**

Technician passes their ID in the request body. In the future, this will come from the JWT.

```go
func (h *TicketHandler) ClaimTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid ticket ID format",
		})
	}

	// For mock auth: technician_id comes from the body.
	// In production this would come from the JWT.
	var body struct {
		TechnicianID string `json:"technician_id" validate:"required,uuid"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := h.service.ClaimTicket(c.Context(), id, body.TechnicianID); err != nil {
		if errors.Is(err, service.ErrClaimConflict) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"error":   "Ticket is already claimed or not available",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Ticket claimed"})
}
```

- [ ] **Step 2: Add transition handlers**

Each handler calls the corresponding service method. They all follow the same pattern.

```go
func (h *TicketHandler) CompleteDiagnosis(c *fiber.Ctx) error {
	return h.handleTransition(c, h.service.CompleteDiagnosis)
}

func (h *TicketHandler) ApproveRepair(c *fiber.Ctx) error {
	return h.handleTransition(c, h.service.ApproveRepair)
}

func (h *TicketHandler) CancelRepair(c *fiber.Ctx) error {
	return h.handleTransition(c, h.service.CancelRepair)
}

func (h *TicketHandler) CompleteRepair(c *fiber.Ctx) error {
	return h.handleTransition(c, h.service.CompleteRepair)
}

func (h *TicketHandler) MarkPickedUp(c *fiber.Ctx) error {
	return h.handleTransition(c, h.service.MarkPickedUp)
}

// handleTransition is a DRY helper for all status transition endpoints.
func (h *TicketHandler) handleTransition(c *fiber.Ctx, action func(ctx context.Context, id string) error) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid ticket ID format",
		})
	}

	if err := action(c.Context(), id); err != nil {
		if errors.Is(err, service.ErrTicketNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Ticket not found",
			})
		}
		if errors.Is(err, service.ErrInvalidTransition) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	return c.JSON(fiber.Map{"success": true})
}
```

- [ ] **Step 3: Add GetBoard handler**

```go
func (h *TicketHandler) GetBoard(c *fiber.Ctx) error {
	board, err := h.service.ListForBoard(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    board,
	})
}
```

- [ ] **Step 4: Add `"context"` to imports**

- [ ] **Step 5: Commit**

```bash
git add apps/backend/internal/handler/
git commit -m "feat: add ticket lifecycle and board endpoints to handler"
```

---

### Task 6: Wire Routes in main.go

**Files:**
- Modify: `apps/backend/main.go`

- [ ] **Step 1: Add middleware import and register RoleMiddleware**

```go
// Add to imports:
"github.com/denden-dr/openbench/apps/backend/internal/handler/middleware"
```

```go
// Add after recover.New() and before routes:
app.Use(middleware.RoleMiddleware)
```

- [ ] **Step 2: Register new ticket routes**

Replace the existing ticket routes block with:

```go
tickets := api.Group("/tickets")

// Public
tickets.Get("/board", ticketHandler.GetBoard)
tickets.Post("/", ticketHandler.Create)

// User actions
tickets.Post("/:id/approve", ticketHandler.ApproveRepair)
tickets.Post("/:id/cancel", ticketHandler.CancelRepair)

// Technician actions
tickets.Post("/:id/claim", middleware.RequireTechnician, ticketHandler.ClaimTicket)
tickets.Post("/:id/complete-diagnosis", middleware.RequireTechnician, ticketHandler.CompleteDiagnosis)
tickets.Post("/:id/complete-repair", middleware.RequireTechnician, ticketHandler.CompleteRepair)
tickets.Post("/:id/pickup", middleware.RequireTechnician, ticketHandler.MarkPickedUp)

// Detail (keep existing)
tickets.Get("/:id", ticketHandler.GetByID)
```

> **Note:** `GET /board` is registered **before** `GET /:id` so Fiber doesn't treat `board` as a UUID param.

- [ ] **Step 3: Commit**

```bash
git add apps/backend/main.go
git commit -m "feat: wire lifecycle routes and role middleware in main.go"
```

---

### Task 7: End-to-End Verification

- [ ] **Step 1: Start the server**

Run: `cd apps/backend && go run main.go`

- [ ] **Step 2: Create a ticket (User)**

```bash
curl -s -X POST http://localhost:3000/api/v1/tickets \
  -H "Content-Type: application/json" \
  -d '{"device_type":"android","brand":"Samsung","model":"Galaxy S24","issue_description":"Screen cracked","diagnosis_fee":"50000"}' | jq
```

Expected: `201 Created`, status is `received`. Save the `id`.

- [ ] **Step 3: Claim ticket (Technician)**

```bash
curl -s -X POST http://localhost:3000/api/v1/tickets/<TICKET_ID>/claim \
  -H "X-Mock-Role: technician" \
  -H "Content-Type: application/json" \
  -d '{"technician_id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"}' | jq
```

Expected: `200 OK`, `{"success": true, "message": "Ticket claimed"}`.

- [ ] **Step 4: Complete diagnosis (Technician)**

```bash
curl -s -X POST http://localhost:3000/api/v1/tickets/<TICKET_ID>/complete-diagnosis \
  -H "X-Mock-Role: technician" | jq
```

Expected: `200 OK`.

- [ ] **Step 5: Approve repair (User)**

```bash
curl -s -X POST http://localhost:3000/api/v1/tickets/<TICKET_ID>/approve | jq
```

Expected: `200 OK`.

- [ ] **Step 6: Complete repair (Technician)**

```bash
curl -s -X POST http://localhost:3000/api/v1/tickets/<TICKET_ID>/complete-repair \
  -H "X-Mock-Role: technician" | jq
```

Expected: `200 OK`.

- [ ] **Step 7: Mark picked up (Technician)**

```bash
curl -s -X POST http://localhost:3000/api/v1/tickets/<TICKET_ID>/pickup \
  -H "X-Mock-Role: technician" | jq
```

Expected: `200 OK`.

- [ ] **Step 8: Verify board**

```bash
curl -s http://localhost:3000/api/v1/tickets/board | jq
```

Expected: Ticket shows status `picked_up`. Response does NOT contain `issue_description` or `diagnosis_fee`.

- [ ] **Step 9: Test cancellation flow**

Create a second ticket, claim it, complete diagnosis, then cancel:

```bash
curl -s -X POST http://localhost:3000/api/v1/tickets/<TICKET_2_ID>/cancel | jq
```

Expected: `200 OK`. Ticket status is `cancelled`.

- [ ] **Step 10: Verify forbidden access**

Attempt to claim without technician role:

```bash
curl -s -X POST http://localhost:3000/api/v1/tickets/<TICKET_ID>/claim \
  -H "Content-Type: application/json" \
  -d '{"technician_id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"}' | jq
```

Expected: `403 Forbidden`.

- [ ] **Step 11: Final commit**

```bash
git commit --allow-empty -m "test: verify full business cycle end-to-end"
```
