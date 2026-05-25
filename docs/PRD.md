# PRD - Phone Repair Management System

## 1. Project Overview
A single-user admin dashboard designed for local phone repair business owners to track and manage repair tickets in one central web interface. The system replaces complex multi-role workflows and external tracking boards with a direct, single-user dashboard that contains all repair statuses, intake forms, ticket search, and quick stat summary metrics.

## 2. User Roles & Permissions

| Role | Access Level | Primary Responsibilities |
| :--- | :--- | :--- |
| **Admin/Owner** | Full Access | Full dashboard access: intake new devices, update repair statuses, set payment outcomes, search tickets, and track repair metrics. |

*Note: Authentication and multi-user roles are removed to simplify local management.*

## 3. Functional Requirements

### 3.1 Repair Ticket Intake & Form
*   **Ticket Creation Form:** Admin can register a new device repair by providing:
    *   **Customer Details:** Customer Name and Gender (Male / Female).
    *   **Device Specs:** Brand and Model.
    *   **Issue:** Primary repair issue (e.g., LCD, battery replacement).
    *   **Additional Details:** Accessory checklist (SIM, Case, SD Card, etc.) and custom notes.
    *   **Financials:** Quote price and Warranty period in days (defaults to 30).
*   **Struct Validation:** Inputs like Gender, Status, and Payment Status are fully validated on submission. Warranty days default to 30 if not provided or <= 0.
*   **Idempotency:** All mutation requests (`POST`/`PATCH`) include an `X-Idempotency-Key` header (UUID v4). Duplicate keys with the same payload return cached `201`; duplicate keys with different payloads return `409 Conflict`.

### 3.2 Main Repair Admin Dashboard
*   **Performance Metrics (Stats Cards):** Real-time calculations of:
    *   **Total Revenue:** Total price of completed and paid repairs.
    *   **Active Repairs:** Count of tickets in active progress.
    *   **Completed Today:** Number of tickets completed on the current date.
    *   **Unpaid Repairs:** Value/count of repairs ready or completed but not yet paid.
*   **Ticket Management Table:** List showing all registered tickets with filters/sorting.
*   **Inline Edit Drawer:** Select any ticket to slide out a quick-edit panel for updating statuses, payment states, warranties, and pricing.
*   **Ticket Deletion:** Admin can permanently delete tickets from the registry.
*   **Mock API Mode:** Frontend can run independently via `MOCK_API=true` with fake data handlers for development.

### 3.3 Status & Workflows
*   **Ticket Statuses (canonical enum):**
    *   `service_in`: Newly received ticket (default on intake).
    *   `on_process`: Device is actively being diagnosed or repaired.
    *   `waiting_confirmation`: Repair done, awaiting customer decision.
    *   `cancelled`: Customer declined or cannot proceed.
    *   `fixed`: Repair complete, waiting for customer collection.
    *   `picked_up`: Device retrieved by customer. This is the payment moment.
*   **Payment Statuses:**
    *   `unpaid`: Repair is not yet settled.
    *   `paid`: Payment successfully received.
*   **Lifecycle Invariants:**
    *   When moving to `picked_up`:
        *   `exit_date` is auto-recorded (UTC now).
        *   `payment_status` auto-set to `paid` (if not explicitly provided).
        *   `warranty_expiry_date` is computed dynamically as `exit_date + warranty_days` — **not stored in the database**.
    *   Moving from `picked_up` back to any other status clears `exit_date`.
    *   Validation rules enforced by backend:
        *   `picked_up` tickets must have `payment_status=paid` and a non-nil `exit_date`.
        *   Non-`picked_up` tickets must have a nil `exit_date`.
        *   Price cannot be negative.
        *   Warranty days cannot be negative.
*   **Cancellation:**
    *   `cancelled` is a valid ticket status (not a deletion). Tickets can transition to `cancelled` through the edit drawer with a technician issue reason.
    *   Deletion permanently removes the ticket from the registry.

## 4. Technical Architecture

### 4.1 Frontend (Svelte)
*   **Framework:** SvelteKit (SPA, adapter-node for production builds).
*   **Reactivity:** Svelte 5 Runes mode (`$state`, `$derived`, `$effect`).
*   **Styling:** Tailwind CSS v4 with CSS variables.
*   **Component Structure:** Main dashboard layout (`+page.svelte`), intake modal, slide-out drawer, stats card, and search filter modules.
*   **Mock Layer:** MSW-inspired manual mock handlers under `src/lib/mocks/` activated via `MOCK_API=true`.
*   **Icons:** Lucide Svelte for UI iconography.

### 4.2 Backend (Go + Fiber)
*   **API Framework:** Fiber v2 (Go HTTP web framework).
*   **Database:** PostgreSQL 16 via `sqlx` for type-safe queries.
*   **Layered Architecture:** handler → service → repository (strict separation).
*   **Idempotency:** Postgres-backed idempotency middleware (`database/idempotency_storage.go`). Stores request hashes keyed by `X-Idempotency-Key`. In-memory by default (single-instance safe); switch to DB-backed for multi-instance deployments.
*   **DTOs:** Request/response types in `internal/dto/` — `CreateTicketRequest`, `UpdateTicketRequest`, `TicketResponse`.
*   **Error Handling:** Centralized `ErrorHandler` middleware. Domain errors (`AppError`) with typed codes propagate from service → handler → HTTP response.
*   **CORS & Proxy:** Vite dev server proxies `/api/*` to the backend (`localhost:3000`).

### 4.3 Testing
*   **Unit Tests:** Go standard `testing` + `testify` assertions. Test every layer (model, service, repository, handler, middleware, config).
*   **Integration Tests:** Build-tagged (`//go:build integration`) tests using Testcontainers for PostgreSQL. Cover handler routes, database queries, and idempotency middleware against a real DB instance.
*   **Mocks:** Testify mocks generated via Mockery for `repository.TicketRepository` and `database.IdempotencyStorage`.
*   **Frontend Checks:** `svelte-check` for TypeScript + Svelte type safety; `npm run build` for production validation.

## 5. Data Model (Schema)

### 5.1 Tickets Table
All details are stored in a single table for maximum query simplicity.

*   `id`: UUID (Primary Key, auto-generated)
*   `customer_name`: Text (Name of the customer)
*   `customer_gender`: Text (Validated: Male, Female)
*   `brand`: Text (Device brand)
*   `model`: Text (Device model)
*   `issue`: Text (Primary repair request details)
*   `additional_description`: Text (Optional, detailed repair notes)
*   `accessories`: Text (List of physical accessories left with the device)
*   `price`: Decimal (Total cost of repair)
*   `status`: Text (Validated: `service_in`, `on_process`, `waiting_confirmation`, `cancelled`, `fixed`, `picked_up`)
*   `payment_status`: Text (Validated: `unpaid`, `paid`)
*   `warranty_days`: Integer (Warranty coverage duration, default 30)
*   `entry_date`: Timestamp (Date device entered the shop)
*   `exit_date`: Timestamp (Date device left the shop or was completed; nil until `picked_up`)

**Computed field** (not stored):
*   `warranty_expiry_date`: Calculated as `exit_date + warranty_days` via the `Ticket.WarrantyExpiryDate()` method. Returns nil if `exit_date` is nil.

---
*End of PRD*
