# Backend Warranty Claims Implementation Plan (Decoupled Queue)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create the database schema, domain model, repository, service, handler, and routing layers in Go to support a decoupled queue model: Intake registers the claim first as `waiting_inspection`, and a subsequent Inspection decision either approves it (spawning a new Rp 0 ticket in a transaction) or voids it.

**Architecture:** 
- Table `warranty_claims` stores intake info (issue, optional description) and is linked to the original ticket.
- `claim_ticket_id` is nullable and populated only when the claim is approved.
- Decoupled API endpoints manage intake creation, queue listing, and approval/void decisions.

**Tech Stack:** Go 1.25, Fiber (v2), sqlx, PostgreSQL.

---

## Files to Create/Modify
- `apps/backend/migrations/000007_create_warranty_claims_table.up.sql` (Create)
- `apps/backend/migrations/000007_create_warranty_claims_table.down.sql` (Create)
- `apps/backend/internal/model/ticket.go` (Modify: Add fields for warranty relationship)
- `apps/backend/internal/model/warranty_claim.go` (Create: Domain model for claims)
- `apps/backend/internal/repository/ticket_repo.go` (Modify: Add columns to SELECT and INSERT/UPDATE queries)
- `apps/backend/internal/repository/warranty_claim_repo.go` (Create: sqlx repository for claims)
- `apps/backend/internal/dto/ticket_dto.go` (Modify: Include new fields in response and request)
- `apps/backend/internal/dto/warranty_claim_dto.go` (Create: Request and response structures)
- `apps/backend/internal/service/ticket_service.go` (Modify: Map the new fields in response mapper)
- `apps/backend/internal/service/warranty_claim_service.go` (Create: Service layer with transaction safety)
- `apps/backend/internal/handler/warranty_claim_handler.go` (Create: Fiber handlers for endpoints)
- `apps/backend/main.go` (Modify: Inject dependencies and register route endpoints)
- `apps/backend/internal/service/warranty_claim_service_test.go` (Create: Unit tests for claims validation logic)

---

### Task 1: Database Migration & Schema Setup

**Files:**
- Create: `apps/backend/migrations/000007_create_warranty_claims_table.up.sql`
- Create: `apps/backend/migrations/000007_create_warranty_claims_table.down.sql`
- Modify: `apps/backend/internal/model/ticket.go`

- [ ] **Step 1: Create the SQL migration up script**

Create `/home/denden/Documents/denden/personal/personal-project/OpenBench/apps/backend/migrations/000007_create_warranty_claims_table.up.sql` with content:
```sql
ALTER TABLE tickets ADD COLUMN is_warranty BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE tickets ADD COLUMN parent_ticket_id UUID REFERENCES tickets(id) ON DELETE SET NULL;

CREATE TABLE warranty_claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    claim_ticket_id UUID REFERENCES tickets(id) ON DELETE SET NULL,
    issue TEXT NOT NULL,
    additional_description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'waiting_inspection',
    void_reason TEXT,
    inspected_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_warranty_claims_ticket_id ON warranty_claims(ticket_id);
CREATE INDEX idx_warranty_claims_status ON warranty_claims(status);
```

- [ ] **Step 2: Create the SQL migration down script**

Create `/home/denden/Documents/denden/personal/personal-project/OpenBench/apps/backend/migrations/000007_create_warranty_claims_table.down.sql` with content:
```sql
DROP TABLE IF EXISTS warranty_claims;
ALTER TABLE tickets DROP COLUMN IF EXISTS parent_ticket_id;
ALTER TABLE tickets DROP COLUMN IF EXISTS is_warranty;
```

- [ ] **Step 3: Modify ticket.go model to add new fields**

Update `apps/backend/internal/model/ticket.go` to include fields:
```go
type Ticket struct {
	ID                    string              `db:"id" json:"id"`
	CustomerName          string              `db:"customer_name" json:"customer_name"`
	CustomerGender        string              `db:"customer_gender" json:"customer_gender"`
	Brand                 string              `db:"brand" json:"brand"`
	Model                 string              `db:"model" json:"model"`
	Issue                 string              `db:"issue" json:"issue"`
	AdditionalDescription *string             `db:"additional_description" json:"additional_description"`
	Accessories           *string             `db:"accessories" json:"accessories"`
	Price                 decimal.Decimal     `db:"price" json:"price"`
	Status                TicketStatus        `db:"status" json:"status"`
	PaymentStatus         TicketPaymentStatus `db:"payment_status" json:"payment_status"`
	WarrantyDays          int                 `db:"warranty_days" json:"warranty_days"`
	EntryDate             time.Time           `db:"entry_date" json:"entry_date"`
	ExitDate              *time.Time          `db:"exit_date" json:"exit_date"`
	IsWarranty            bool                `db:"is_warranty" json:"is_warranty"`
	ParentTicketID        *string             `db:"parent_ticket_id" json:"parent_ticket_id"`
}
```

- [ ] **Step 4: Run the migrations locally**

Run: `make migrate-up`
Expected: Migrations apply cleanly to the local PostgreSQL database.

---

### Task 2: Implement Domain Model & Repository

**Files:**
- Create: `apps/backend/internal/model/warranty_claim.go`
- Modify: `apps/backend/internal/repository/ticket_repo.go`
- Create: `apps/backend/internal/repository/warranty_claim_repo.go`

- [ ] **Step 1: Create the warranty claim model**

Create `/home/denden/Documents/denden/personal/personal-project/OpenBench/apps/backend/internal/model/warranty_claim.go` with content:
```go
package model

import "time"

type WarrantyClaimStatus string

const (
	ClaimWaitingInspection WarrantyClaimStatus = "waiting_inspection"
	ClaimApproved          WarrantyClaimStatus = "approved"
	ClaimVoid              WarrantyClaimStatus = "void"
)

type WarrantyClaim struct {
	ID                    string              `db:"id" json:"id"`
	TicketID              string              `db:"ticket_id" json:"ticket_id"`
	ClaimTicketID         *string             `db:"claim_ticket_id" json:"claim_ticket_id"`
	Issue                 string              `db:"issue" json:"issue"`
	AdditionalDescription *string             `db:"additional_description" json:"additional_description"`
	Status                WarrantyClaimStatus `db:"status" json:"status"`
	VoidReason            *string             `db:"void_reason" json:"void_reason"`
	InspectedAt           *time.Time          `db:"inspected_at" json:"inspected_at"`
	CreatedAt             time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time           `db:"updated_at" json:"updated_at"`
}
```

- [ ] **Step 2: Update existing ticket_repo.go to support new columns**

Modify `apps/backend/internal/repository/ticket_repo.go` to select, insert, and update `is_warranty` and `parent_ticket_id`.
Ensure `Create` and `CreateTx` methods correctly support these columns.

- [ ] **Step 3: Create the warranty claim repository interface and implementation**

Create `/home/denden/Documents/denden/personal/personal-project/OpenBench/apps/backend/internal/repository/warranty_claim_repo.go`. It should support:
- `Create(ctx context.Context, claim *model.WarrantyClaim) error`
- `GetByID(ctx context.Context, id string) (*model.WarrantyClaim, error)`
- `List(ctx context.Context, status string) ([]*model.WarrantyClaim, error)`
- `UpdateTx(ctx context.Context, tx *sqlx.Tx, claim *model.WarrantyClaim) error`
- `CreateTicketTx(ctx context.Context, tx *sqlx.Tx, ticket *model.Ticket) error`

---

### Task 3: DTOs & Service Layer

**Files:**
- Modify: `apps/backend/internal/dto/ticket_dto.go`
- Create: `apps/backend/internal/dto/warranty_claim_dto.go`
- Modify: `apps/backend/internal/service/ticket_service.go`
- Create: `apps/backend/internal/service/warranty_claim_service.go`

- [ ] **Step 1: Update Ticket DTO structs**

Modify `apps/backend/internal/dto/ticket_dto.go` to include `IsWarranty` and `ParentTicketID`.

- [ ] **Step 2: Create Warranty Claim DTOs**

Create `/home/denden/Documents/denden/personal/personal-project/OpenBench/apps/backend/internal/dto/warranty_claim_dto.go` supporting:
- `CreateWarrantyClaimRequest` (contains `TicketID`, `Issue`, `AdditionalDescription`)
- `VoidWarrantyClaimRequest` (contains `VoidReason`)
- `WarrantyClaimResponse`
- `ClaimCreationResult` (composition of claim + ticket)

- [ ] **Step 3: Create the Warranty Claim service**

Create `/home/denden/Documents/denden/personal/personal-project/OpenBench/apps/backend/internal/service/warranty_claim_service.go` containing:
- `CreateClaim`: Validates original ticket, checks expiry, and registers claim as `waiting_inspection`.
- `ListClaims`: Returns filtered/unfiltered warranty claims queue.
- `ApproveClaim`: Uses `db.BeginTxx` to atomically:
  1. Retrieve claim.
  2. Create a new `Ticket` with price 0, `is_warranty = true`, and status `on_process`.
  3. Update claim status to `approved` and save the `claim_ticket_id`.
- `VoidClaim`: Updates claim status to `void` and records the `void_reason`.

---

### Task 4: Handler & Endpoint Routing Setup

**Files:**
- Create: `apps/backend/internal/handler/warranty_claim_handler.go`
- Modify: `apps/backend/main.go`

- [ ] **Step 1: Create HTTP Handler**

Create `/home/denden/Documents/denden/personal/personal-project/OpenBench/apps/backend/internal/handler/warranty_claim_handler.go` with functions:
- `Create(c *fiber.Ctx)`
- `List(c *fiber.Ctx)`
- `Approve(c *fiber.Ctx)`
- `Void(c *fiber.Ctx)`

- [ ] **Step 2: Inject Dependencies and Register endpoints in main.go**

Under the existing DI setup, wire `WarrantyClaimRepository`, `WarrantyClaimService`, and `WarrantyClaimHandler`.
Register routes:
- `POST /api/v1/warranty-claims`
- `GET /api/v1/warranty-claims`
- `POST /api/v1/warranty-claims/:id/approve`
- `POST /api/v1/warranty-claims/:id/void`

---

### Task 5: Implement Service Unit Tests & Verification

**Files:**
- Create: `apps/backend/internal/service/warranty_claim_service_test.go`

- [ ] **Step 1: Write warranty claim service validation unit tests**
- [ ] **Step 2: Run Go unit tests and verify they pass (`make test-backend-unit`)**
- [ ] **Step 3: Run database integration tests if applicable (`make test-backend-integration`)**
