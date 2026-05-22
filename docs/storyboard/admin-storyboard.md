# Admin Storyboard

This document focuses on the workflows from the perspective of the Admin managing the back-office operations.

## Admin Stories

### Story 10: Admin — Inventory & User Management

> Admin manages the back-office.

1. Admin logs in → sees dashboard (revenue, active tickets, low stock alerts)
2. Admin navigates to Inventory Management
3. Admin adds new parts: name, brand compatibility, grade (ODM/Original), price, stock
4. Admin adjusts stock levels → audit log created
5. Admin navigates to User Management
6. Admin creates technician account
7. Admin navigates to Payment Management
8. Admin views pending/completed payments, reconciles cash payments
9. Admin navigates to Audit Log
10. Admin reviews sensitive actions: status changes, stock deductions, price updates

---

### Admin Roles in Customer & Ticket Flows

The Admin is also responsible for handling specific ticket overrides and cancellations:

#### Cancelling Tickets (from Story 2, 3, & 4)
- **Cancellation** can only be performed by **Admin** per domain rules (not directly by customer or technician).
- Admin cancels ticket with reason logged `[cancelled]` when:
  - Technician determines device is unrepairable.
  - Customer declines to proceed after initial diagnosis.
  - Customer declines to continue after a mid-repair re-diagnosis.

#### Restocking Parts (from Story 6)
- When a technician flags a ticket as `[waiting_parts]`, the Admin is responsible for procuring and restocking the required parts (inventory updated).
