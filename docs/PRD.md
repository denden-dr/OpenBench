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
    *   **Financials:** Quote price and Warranty period in days.
*   **Struct Validation:** Inputs like Gender, Status, and Payment Status are fully validated on submission.

### 3.2 Main Repair Admin Dashboard
*   **Performance Metrics (Stats Cards):** Real-time calculations of:
    *   **Total Revenue:** Total price of completed and paid repairs.
    *   **Active Repairs:** Count of tickets in active progress.
    *   **Completed Today:** Number of tickets completed on the current date.
    *   **Unpaid Repairs:** Value/count of repairs ready or completed but not yet paid.
*   **Ticket Management Table:** List showing all registered tickets with filters/sorting.
*   **Inline Edit Drawer:** Select any ticket to slide out a quick-edit panel for updating statuses, payment states, warranties, and pricing.
*   **Ticket Deletion:** Admin can permanently delete tickets from the registry.

### 3.3 Status & Workflows
*   **Repair Statuses:**
    *   `service_in`: Newly received ticket.
    *   `picked_up`: Ticket claimed or in active diagnosis/repair.
    *   `done`: Repair successfully completed.
    *   `cancel`: Repair cancelled or unrepairable.
*   **Payment Statuses:**
    *   `unpaid`: Repair is not yet settled.
    *   `paid`: Payment successfully received.
*   **Dates Logic:**
    *   When moving to `picked_up`, the entry date is recorded.
    *   When moving to `done` or `cancel`, the exit date is logged and the warranty expiry is automatically calculated based on the ticket's warranty days.

## 4. Technical Architecture

### 4.1 Frontend (Svelte)
*   **Framework:** SvelteKit (Single Page App style).
*   **Styling:** Vanilla CSS styled with custom CSS variables and Tailwind classes.
*   **Component Structure:** Main dashboard layout (`+page.svelte`), intake modal, slide-out drawer, stats card, and search filter modules.

### 4.2 Backend (Go + Fiber)
*   **API Framework:** Fiber (Go HTTP web framework).
*   **Database:** PostgreSQL (Relational storage).
*   **CORS & Proxy:** Local Vite development proxy configs to forward frontend `/api/*` traffic to the backend server.

## 5. Data Model (Schema)

### 5.1 Tickets Table
All details are stored in a single table for maximum query simplicity.

*   `id`: UUID (Primary Key, automatically generated)
*   `customer_name`: Text (Name of the customer)
*   `customer_gender`: Text (Validated: Male, Female)
*   `brand`: Text (Device brand)
*   `model`: Text (Device model)
*   `issue`: Text (Primary repair request details)
*   `additional_description`: Text (Optional, detailed description or repair notes)
*   `accessories`: Text (List of physical accessories left with the device)
*   `price`: Decimal (Total cost of repair)
*   `status`: Text (Validated: `service_in`, `picked_up`, `done`, `cancel`)
*   `payment_status`: Text (Validated: `unpaid`, `paid`)
*   `warranty_days`: Integer (Warranty coverage duration)
*   `entry_date`: Timestamp (Date device entered the shop)
*   `exit_date`: Timestamp (Date device left the shop or was completed)
*   `warranty_expiry_date`: Timestamp (Automatically calculated upon completion)

---
*End of PRD*
