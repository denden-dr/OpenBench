# Future Implementation Plan

Based on the PRD (`docs/PRD.md`) and the current project snapshot (`docs/project-snapshot.md`), here is the checklist of domains, features, and systems that still need to be implemented.

### 1. Authentication & User Management
The current backend setup lacks real user identity management.
- [ ] **Database Setup**: Implement `Users` and `Customers` tables.
- [ ] **Supabase Integration**: Integrate Supabase Auth for JWT verification in Go Fiber.
- [ ] **Role-based Access Control (RBAC)**: Replace the mocked role middleware with real JWT parsing to differentiate between `Guest`, `Customer`, `Technician`, and `Admin`.

### 2. Ticket Domain Enhancements
While the basic ticket lifecycle is in place, it lacks several specific PRD requirements.
- [ ] **Booking Info Expansion**: Add `device_access_info` (passcode/pattern) and `accessories` tracking to `CreateTicketRequest` and database schema.
- [ ] **Estimations**: Require technicians to set an `Estimated Completion Time` when submitting a diagnosis.
- [ ] **PDF Receipts**: Implement receipt generation with a QR code for intake proof.
- [ ] **Internal Comments**: Implement the `TicketComments` table and endpoints for staff to leave private technical notes.

### 3. Missing Backend Domains
Entire domains defined in the PRD have not been scaffolded yet.
- [ ] **Inventory & Parts Management**: Implement `Parts` and `TicketParts` tables, allowing technicians to log OEM vs. Original parts and track stock levels.
- [ ] **Warranty Management**: Implement logic to automatically calculate warranty expiry (7 days vs 30 days) based on the grade of parts used upon ticket completion.
- [ ] **Payments & Invoicing**: Implement the `Payments` table and integrate the Midtrans payment gateway to auto-calculate `Diagnosis Fee + Labor Fee + Parts`.
- [ ] **Media Documentation**: Implement the `Attachments` table and integrate Supabase Storage / S3 to allow technicians to upload "Before/After" diagnostic photos.
- [ ] **Audit Logging**: Create the `AuditLogs` table to automatically track sensitive actions (status changes, stock deductions).

### 4. Frontend Application (SvelteKit)
The `apps/frontend` directory is currently just a scaffold.
- [ ] **Auth Setup**: Configure Supabase SSR client for secure, cookie-based authentication.
- [ ] **Customer UI**: Build the multi-step booking form, public tracker (search by ID + phone number), and user profile page.
- [ ] **Technician UI**: Build the Kanban-style ticket queue, ticket claim interface, and status update modals.
- [ ] **Admin UI**: Build the dashboard for managing inventory, viewing system audits, and tracking revenue reports.
