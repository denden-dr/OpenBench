# System Architecture - PhoneFix

This document outlines the technical flow and interaction between the Client, Frontend, Backend, and Third-party services.

## 1. High-Level Flow (Authentication & Data)

The following sequence diagram illustrates how a user request is handled from the browser through to the database, using Supabase for identity management.

```mermaid
sequenceDiagram
    participant Client as Browser (Client)
    participant Svelte as SvelteKit (Frontend)
    participant Supabase as Supabase Auth
    participant Fiber as Go Fiber (Backend API)
    participant DB as PostgreSQL

    Note over Client, Supabase: Authentication Phase
    Client->>Svelte: Login / Register (POST)
    Svelte->>Supabase: Verify Credentials
    Supabase-->>Svelte: Return Session & JWT
    Svelte-->>Client: Set-Cookie: sb-access-token (HttpOnly, Secure)

    Note over Client, DB: Authorized Request Phase
    Client->>Fiber: API Call (Cookies sent automatically)
    
    Fiber->>Fiber: Extract JWT from Cookie
    Fiber->>Supabase: Verify JWT
    Supabase-->>Fiber: Token Valid (User ID: 123)

    Fiber->>DB: SQL Operations (Insert Ticket)
    DB-->>Fiber: Success
    
    Fiber-->>Svelte: JSON Response (201 Created)
    Svelte-->>Client: Update UI (Ticket #456 Created)
```

## 2. Specific Feature Flows

### 2.1 Ticket Claiming Flow (First-Come, First-Served)
This flow handles how a technician "Takes" a ticket from the public queue. It uses a conditional update to prevent double-claiming.

```mermaid
sequenceDiagram
    participant Tech as Technician (UI)
    participant API as Go Fiber API
    participant DB as PostgreSQL

    Tech->>API: POST /tickets/{id}/claim
    
    API->>DB: UPDATE tickets SET technician_id = {TechID}, status = 'diagnosing' WHERE id = {id} AND technician_id IS NULL
    
    alt Claim Successful (Rows Affected > 0)
        DB-->>API: Success
        API-->>Tech: 200 OK (Ticket Assigned)
    else Already Claimed (Rows Affected = 0)
        DB-->>API: No change
        API-->>Tech: 409 Conflict (Ticket already taken by another tech)
    end
```

### 2.2 Repair & Inventory Workflow
This flow ensures that when a technician marks a repair as "In-Progress" or "Completed", the parts used are correctly deducted from the inventory within a database transaction.

```mermaid
sequenceDiagram
    participant Tech as Technician (UI)
    participant API as Go Fiber API
    participant DB as PostgreSQL

    Tech->>API: POST /tickets/{id}/parts (PartID, Qty)
    
    Note over API, DB: Database Transaction Start
    API->>DB: SELECT stock_level FROM parts WHERE id = {PartID}
    DB-->>API: stock_level: 10
    
    alt Stock Available
        API->>DB: UPDATE parts SET stock_level = stock_level - {Qty}
        API->>DB: INSERT INTO ticket_parts (ticket_id, part_id, quantity)
        API->>DB: UPDATE tickets SET status = 'repairing'
        Note over API, DB: Transaction Commit
        API-->>Tech: 200 OK (Status Updated & Parts Logged)
    else Out of Stock
        Note over API, DB: Transaction Rollback
        API-->>Tech: 400 Bad Request (Insufficient Stock)
    end
```

### 2.3 Online Payment Flow (Webhooks)
This flow handles the asynchronous update of payment status when a customer pays via a gateway like Midtrans.

```mermaid
sequenceDiagram
    participant Client as Customer
    participant PGW as Payment Gateway (Midtrans)
    participant API as Go Fiber API
    participant DB as PostgreSQL

    Client->>PGW: Complete Payment (QRIS/Bank Transfer)
    PGW-->>Client: Payment Success
    
    Note over PGW, API: Asynchronous Webhook
    PGW->>API: POST /webhooks/payment (Signature, OrderID, Status)
    
    API->>API: Verify Webhook Signature
    API->>DB: UPDATE payments SET status = 'completed' WHERE external_id = {OrderID}
    API->>DB: UPDATE tickets SET status = 'completed' WHERE id = {TicketID}
    
    API-->>PGW: 200 OK
```

## 3. Component Responsibilities

### 3.1 SvelteKit (Frontend)
*   Handles routing and page rendering.
*   **Session Management**: Uses `@supabase/ssr` to manage auth sessions in `HttpOnly` cookies.
*   **Auth Proxy/Server Actions**: Handles login/logout server-side to set/clear cookies.
*   Proxies/Calls the Go Backend; browser attaches cookies automatically.
*   Provides a premium, responsive UI with vanilla CSS.

### 3.2 Go Fiber (Backend)
*   Exposes a RESTful API.
*   **Auth Middleware**: Extracts the JWT from the `Cookie` header (e.g., `sb-access-token`) and verifies it with Supabase.
*   **Business Logic**: Handles complex validations, invoicing logic, and inventory math.
*   **Database Access**: Communicates with PostgreSQL.

### 3.3 Supabase Auth
*   Identity Provider (IdP).
*   Handles email/password and social login providers.
*   Issues short-lived JWTs for secure API access.

### 3.4 PostgreSQL
*   The source of truth for business data (Tickets, Inventory, Customer details).
*   Linked to Supabase users via a `supabase_uid` mapping.

---
*Last Updated: 2026-04-30*
