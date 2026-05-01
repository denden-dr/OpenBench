# PRD - Phone Repair Management System

## 1. Project Overview
A web application designed to manage the end-to-end workflow of a phone repair business. The system serves both customers (booking and tracking) and staff (workflow and inventory management).

## 2. User Roles & Permissions

| Role | Access Level | Primary Responsibilities |
| :--- | :--- | :--- |
| **Guest/Public** | Unauthenticated | Track repair status via Ticket ID, view service prices. |
| **Customer** | Authenticated | Manage profile, book new repairs, view full repair history. |
| **Technician** | Authenticated | View assigned tickets, update repair status, log parts used, add technical notes. |
| **Admin** | Authenticated | Full system access: manage inventory, users, financial reporting, and system settings. |

## 3. Functional Requirements

### 3.1 Customer Features
*   **Booking System:** Multi-step form to select Device -> Brand -> Model -> Issue.
    *   **Device Access:** Request **Passcode/Pattern** OR instruction to enable **Repair/Maintenance Mode** (for supported devices).
    *   **Fee Transparency:** Explicitly display the **Mandatory Diagnosis Fee** before booking.
    *   **Accessory Tracking:** Checklist for items left with the device (SIM, Case, SD Card, etc.).
    *   **Terms Agreement:** Required checkbox acknowledging the diagnosis fee and terms of service.
*   **Public Tracker:** Search by Ticket ID + Phone Number to see real-time status.
*   **Intake Receipt:** Ability to download/print a PDF receipt with a **QR Code** for easy status tracking and proof of drop-off.
*   **Profile Management:** View active and past repairs, download invoices.

### 3.2 Staff & Technician Features
*   **Ticket Queue:** View all unassigned tickets that are ready for diagnosis.
*   **Claim Ticket:** Technicians can "Take" a ticket, which assigns it to them and updates the status to 'diagnosing'.
*   **Diagnosis & Estimating:** Add technical notes, diagnostic results, and a mandatory **Estimated Completion Time** (visible to customer).
*   **Media Documentation:** Upload "Before" photos during diagnosis and "After" photos upon completion for liability protection.
*   **Internal Comments:** Add private notes or technical updates that aren't visible to the customer.
*   **Status Management:** One-click status updates for their assigned tickets with return-accessory reminders.

### 3.3 Admin & Operations
*   **Inventory Manager:** Track stock levels and specify **Part Grade** (ODM vs Original) for each item.
*   **Warranty Management:** Automatically calculate warranty expiry based on part grade (ODM: 7 days, Original: 30 days).
*   **Invoicing & Payments:** Auto-calculate totals: `Diagnosis Fee + Labor Fee + Parts`. Support for **Cash** and **Online Payments** (Midtrans).
*   **Audit Logging:** Automatic tracking of all sensitive actions (status changes, price updates, inventory adjustments).
*   **Reporting:** Monthly revenue, most common issues, and technician performance.

## 4. Technical Architecture

### 4.1 Frontend (Svelte)
*   **Framework:** SvelteKit (for routing and SSR/SPA hybrid).
*   **Styling:** Vanilla CSS / Modern CSS Variables (for high-performance, premium feel).
*   **State Management:** Svelte Stores for auth and UI state.

### 4.2 Backend (Go + Fiber)
*   **API Framework:** Fiber (Express-like, high performance).
*   **Database:** PostgreSQL (Relational data for tickets, parts, and users).
*   **Auth:** Supabase Auth integration. Backend will verify Supabase JWTs.
*   **File Storage:** Supabase Storage or S3-compatible storage.
### 4.3 System Flow
For a detailed sequence diagram of the Client -> Frontend -> Backend flow, see **[ARCHITECTURE.md](./ARCHITECTURE.md)**.

## 5. Data Model (Schema)

### 5.1 Users Table
*   `id`: Primary Key (UUID)
*   `supabase_uid`: UUID (Unique, maps to Supabase auth.users.id)
*   `email`: String (Unique)
*   `role`: Enum (admin, technician, customer)
*   `name`: String

### 5.2 Customers Table
*   `id`: Primary Key (UUID)
*   `user_id`: Foreign Key (Users.id)
*   `phone`: String
*   `address`: Text

### 5.3 Tickets Table
*   `id`: Primary Key (UUID/ShortID)
*   `customer_id`: Foreign Key (Customers.id)
*   `device_type`: Enum (Android, Apple)
*   `brand`: String
*   `model`: String
*   `issue_description`: Text
*   `device_access_info`: String (Passcode or "Repair Mode Enabled")
*   `accessories`: JSON/Text (List of items like SIM, Case, etc.)
*   `diagnosis_fee`: Decimal (Mandatory)
*   `labor_fee`: Decimal
*   `status`: Enum (received, diagnosing, waiting_parts, repairing, ready, completed, cancelled)
*   `technician_id`: Foreign Key (Users.id, Nullable)
*   `estimated_ready_at`: Timestamp (Nullable, set during diagnosis)
*   `warranty_expiry`: Timestamp (Nullable, calculated upon completion)
*   `created_at`: Timestamp
*   `updated_at`: Timestamp

### 5.4 Parts Table
*   `id`: Primary Key (UUID)
*   `name`: String
*   `brand_compatibility`: String
*   `grade`: Enum (ODM, Original)
*   `price`: Decimal
*   `stock_level`: Integer

### 5.5 TicketParts (Junction Table)
*   `ticket_id`: Foreign Key (Tickets.id)
*   `part_id`: Foreign Key (Parts.id)
*   `quantity`: Integer

### 5.6 Attachments Table
*   `id`: Primary Key (UUID)
*   `ticket_id`: Foreign Key (Tickets.id)
*   `file_url`: String
*   `type`: Enum (before, after, receipt)
*   `created_at`: Timestamp

### 5.7 Payments Table
*   `id`: Primary Key (UUID)
*   `ticket_id`: Foreign Key (Tickets.id)
*   `amount`: Decimal
*   `method`: Enum (cash, online)
*   `status`: Enum (pending, completed, failed)
*   `external_id`: String (Midtrans/Xendit Order ID)
*   `created_at`: Timestamp

### 5.8 TicketComments Table (Threaded)
*   `id`: Primary Key (UUID)
*   `ticket_id`: Foreign Key (Tickets.id)
*   `user_id`: Foreign Key (Users.id)
*   `content`: Text
*   `created_at`: Timestamp

### 5.9 AuditLogs Table
*   `id`: Primary Key (UUID)
*   `user_id`: Foreign Key (Users.id)
*   `action`: String (e.g., "update_status", "deduct_stock")
*   `entity_type`: String (e.g., "ticket", "part", "payment")
*   `entity_id`: UUID
*   `changes`: JSON (Old vs New values)
*   `created_at`: Timestamp

---
*End of PRD*
