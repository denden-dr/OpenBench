# Frontend UI Pages Spec

## Public

### Landing Page
> Route: `/`

- **Navbar**: Logo, Services & Pricing link, Track Repair link, Public Queue link, Log In / Sign Up button
- **Hero Section**: Tagline + CTA → "Book a Repair" button (navigates to booking form)
- **Service Highlights**: Grid of key selling points (transparent pricing, warranty, photo documentation, real-time tracking)
- **How It Works**: Step-by-step visual (Book → Diagnose → Repair → Pickup)
- **Footer**: Contact info, social links, copyright

### Booking Form Page
> Route: `/book`  
> Auth: Accessible to guests (view & fill), submit requires authentication (deferred sign-up via Story 7a)

- **Device Selection**:
    - Device type: Phone / Tablet (radio / toggle)
    - Brand: Dropdown (dynamic list)
    - Model: Dropdown (filtered by selected brand)
- **Issue Description**:
    - Template issue list (common issues as checkboxes) + "Other" free text
    - Description field: explain when/how the issue occurred — *mandatory*
- **Device Access Info**:
    - Passcode / Pattern / Repair Mode — *mandatory*
- **Accessories Checklist**:
    - SIM card, case, SD card, charger, etc. — *optional*, stored as JSON
- **Diagnosis Fee Acknowledgement**:
    - Fee amount displayed clearly
    - Mandatory checkbox: "I understand and agree to the diagnosis fee"
- **Confirm Button**:
    - If authenticated → submit booking, show Ticket ID + QR Code receipt (PDF download)
    - If guest → save form to localStorage, redirect to sign-up (Story 7a flow)

### Repair Tracking Page
> Route: `/track`  
> Auth: Public (no login required)

- **Search Form**: Ticket ID + Phone Number inputs
- **Result Card** (sanitized):
    - Current status badge (received / diagnosing / repairing / ready / completed / cancelled)
    - Device brand & model
    - Estimated ready date
- **Not Shown**: Internal notes, photos, payment details, customer identity

### Log In / Sign Up Page
> Route: `/login`, `/signup`  
> Auth: Public (redirect to dashboard if already authenticated)

- **Sign Up**:
    - Social login (Google) — *recommended*, shown as primary CTA
    - Email + password — secondary option
    - Fields: Full name, email, phone, password
    - Deferred booking banner: _"Create an account to submit your booking"_ (if redirected from booking form)
- **Log In**:
    - Social login (Google)
    - Email + password
- **Post-auth redirect**:
    - If pending booking in localStorage → redirect to booking form (pre-filled)
    - Otherwise → redirect to profile / ticket list

### Service and Fee Page
> Route: `/services`  
> Auth: Public

- **Service List**: Table or card grid showing:
    - Service type (e.g., Screen Replacement, Battery Swap, Charging Port)
    - Base diagnosis fee
    - Estimated labor fee range
    - Part grade options (ODM vs Original) with price ranges
- **Warranty Info**: Brief explanation of warranty terms (7 days ODM / 30 days Original)

### Public Queue Status Board
> Route: `/queue`  
> Auth: Public

- **Kanban-style board**: Columns by status (Received → Diagnosing → Repairing → Ready)
- **Each card shows**: Masked Ticket ID, brand, model, current status
- **Anonymized**: No names, phone numbers, or identifying info
- **Real-time updates**: Auto-refresh or WebSocket

---

## Registered User

### Profile Page
> Route: `/profile`  
> Auth: Customer

- **Personal Info** (editable):
    - Full name
    - Email (read-only if social auth)
    - Phone number
    - Address
- **Account Actions**:
    - Change password (disabled/hidden if social auth)
    - Log out

### Ticket List Page
> Route: `/tickets`  
> Auth: Customer

- **Active / Completed toggle filter**
- **Ticket Cards**: Each card shows:
    - Ticket ID (short ID)
    - Device: brand + model
    - Current status badge
    - Created date
    - Estimated ready date (if set)
- **Click → navigates to Ticket Details Page**

### Ticket Details Page
> Route: `/tickets/:id`  
> Auth: Customer (owner only)

- **Header**: Ticket ID, status badge, device brand + model
- **Timeline / Status Progress**: Visual step indicator (received → diagnosing → ... → completed)
- **Device Info**: Type, brand, model, issue description, accessories checklist
- **Diagnosis Section** (visible after `diagnosing`):
    - Technician findings
    - Fee breakdown: diagnosis fee + labor fee + parts list with prices
    - Estimated completion date
    - **Action buttons**: Confirm to proceed / Decline (triggers admin cancellation)
- **Repair Updates** (visible during `repairing`):
    - Progress notes (if any public-facing ones exist)
    - Waiting for parts notice (if `waiting_parts`)
- **Before & After Photos** (visible after uploaded):
    - Gallery with privacy controls (blur toggle, download option)
- **Payment Section** (visible at `ready`):
    - Total amount breakdown
    - Payment method selection: Cash / Online (Midtrans)
    - Payment status badge (pending / completed / failed)
    - Retry button (if payment failed)
- **Warranty Info** (visible after `completed`):
    - Warranty expiry date
    - Parts used with grades
- **Invoice Download**: PDF button (available after payment)

---

## Technician User

### Queue List Page
> Route: `/tech/queue`  
> Auth: Technician

- **Tabs / Filters**:
    - Unassigned tickets (claimable queue, sorted by `created_at` oldest first)
    - My active tickets (assigned to current technician)
- **Ticket Cards**: Each card shows:
    - Ticket ID
    - Device: brand + model
    - Issue summary (truncated)
    - Status badge
    - Time since received
- **Claim button** on unassigned tickets (atomic claim, shows 409 toast if already taken)

### Queue Details Page (Ticket Workspace)
> Route: `/tech/tickets/:id`  
> Auth: Technician (assigned or unclaimed)

- **Header**: Ticket ID, status, customer name, device info
- **Customer Info Panel**: Name, phone, issue description, accessories checklist, access info
- **Photo Upload Section**:
    - "Before" photos — upload during `diagnosing`
    - "After" photos — upload during `repairing` / before `ready`
- **Diagnosis Form** (active during `diagnosing`):
    - Technical notes (internal, not visible to customer)
    - Diagnosis findings (visible to customer)
    - Fee breakdown: diagnosis fee, labor fee
    - Estimated completion time (date picker)
    - Mark as unrepairable toggle (sends to admin for cancellation)
    - Submit diagnosis → status moves to `waiting_customer_confirm`
- **Parts Logging (POS-style)** (active during `repairing`):
    - Search/select parts from inventory
    - Set quantity per part
    - Part grade + price auto-filled from inventory
    - Running total displayed
    - Stock validation (insufficient stock → error, blocks submission)
- **Status Actions**:
    - `diagnosing` → Submit diagnosis
    - `waiting_parts` → Flag (with note about needed parts)
    - `repairing` → Mark as ready (triggers accessory return reminder modal)
- **Accessory Return Reminder**: Modal/checklist when marking `ready`
    - SIM card returned? ☐
    - Case returned? ☐
    - SD card returned? ☐
    - Other items from booking? ☐
- **Re-diagnosis Flow** (if critical problem found during `repairing`):
    - "Flag Additional Problem" button → resets to `diagnosing`
    - Updated diagnosis form with original + new findings
    - Submit updated cost estimate → `waiting_customer_confirm`

---

## Administrator

### Dashboard
> Route: `/admin`  
> Auth: Admin

- **KPI Summary Cards**:
    - Total revenue (today / this week / this month)
    - Active tickets count (by status breakdown)
    - Tickets completed today
    - Average repair turnaround time
- **Alerts Panel**:
    - Low stock alerts (parts below threshold)
    - Pending payments requiring reconciliation
    - Tickets in `waiting_parts` status
- **Quick Actions**: Links to Ticket Management, Inventory, Payments

### User Management Page
> Route: `/admin/users`  
> Auth: Admin

- **User Table**: Sortable, searchable list
    - Columns: Name, email, phone, role, status, created date
- **Role Filter**: All / Admin / Technician / Customer
- **Actions**:
    - Create technician / admin account (invite form: name, email, role)
    - Edit user details
    - Deactivate / reactivate user account
- **User Detail Modal/Page**: View user's activity (tickets assigned, login history)

### Ticket Management Page
> Route: `/admin/tickets`  
> Auth: Admin

- **Ticket Table**: Sortable, filterable, searchable
    - Columns: Ticket ID, customer, device, status, technician, created date, updated date
- **Status Filter**: All / per-status dropdown
- **Ticket Detail View** (click to expand or navigate):
    - Full ticket info (same as technician view but read-only for technical fields)
    - **Admin Actions**:
        - Cancel ticket (with mandatory reason field → logged to audit)
        - Reassign technician
        - Override status (emergency use, logged to audit)
        - Confirm cash payment received → moves ticket to `completed`

### Inventory Management Page
> Route: `/admin/inventory`  
> Auth: Admin

- **Parts Table**: Sortable, searchable
    - Columns: Part name, brand compatibility, grade (ODM / Original), price, stock level, last updated
- **Stock Status Indicators**: Color-coded (green = OK, yellow = low, red = out of stock)
- **Actions**:
    - Add new part (form: name, brand compatibility, grade, price, initial stock)
    - Edit part details (price, brand compatibility)
    - Adjust stock level (with reason field → creates audit log entry)
- **Restock Alert Config**: Set low-stock threshold per part

### Payment Management Page
> Route: `/admin/payments`  
> Auth: Admin

- **Payment Table**: Sortable, filterable
    - Columns: Payment ID, Ticket ID, customer name, amount, method (cash / online), status, date
- **Status Filter**: All / Pending / Completed / Failed
- **Actions**:
    - Confirm cash payment → moves ticket to `completed`
    - View Midtrans transaction details (external ID, webhook status)
    - Issue refund (if applicable, logs to audit)
- **Reconciliation View**: Unmatched payments, discrepancies

### Audit Log Page
> Route: `/admin/audit`  
> Auth: Admin

- **Log Table**: Reverse chronological
    - Columns: Timestamp, user (who did it), action type, entity (ticket/part/user), details
- **Filters**:
    - By action type (status change, stock adjustment, price update, payment confirmation, user role change, ticket cancellation)
    - By user
    - By date range
    - By entity type
- **Log Detail**: Expandable row or modal showing full `changes` JSON diff (before → after)

---

## Shared / Cross-Cutting

### Notification System
> Delivery mechanism TBD (email / SMS / in-app)

Referenced in storyboards but not yet a standalone page. Notifications trigger on:
- Diagnosis complete → customer
- Device ready for pickup → customer
- Waiting for parts → customer
- Payment failed → customer
- Low stock alert → admin

### Invoice / Receipt (PDF)
> Not a standalone page — generated as downloadable PDF

- Triggered from Ticket Details Page (customer) or Ticket Management (admin)
- Contains: Ticket ID, QR code, device info, fee breakdown, parts list with grades, warranty info, payment status
- Zero-amount invoice issued for unrepairable devices (record only)

### Error / Empty States
- **404 Page**: Route not found
- **403 Page**: Unauthorized access attempt
- **Empty Ticket List**: "No repairs yet — book your first repair!"
- **No Search Results**: "No ticket found. Double-check your Ticket ID and phone number."
- **Payment Failed**: "Payment unsuccessful. Please try again or choose a different method."