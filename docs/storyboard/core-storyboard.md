# Core Storyboard Index

The core storyboard has been dissected by role into separate files for better maintainability and clarity. Please refer to the specific storyboard based on the actor:

1. [User Storyboard](./user-storyboard.md) — Customer workflows, public tracking, booking, and payment.
2. [Technician Storyboard](./technician-storyboard.md) — Technician daily routines, diagnosis, and repair workflows.
3. [Admin Storyboard](./admin-storyboard.md) — Back-office management, inventory, users, and ticket cancellations.

## General Notes

- **Booking details are the customer's responsibility** — system does not validate device condition claims.
- **Diagnosis fee is waived when unrepairable** — if the technician determines a device cannot be repaired, the customer pays nothing. The diagnosis fee only applies when the customer has a choice (proceed or cancel).
- **`waiting_customer_confirm`** is not in the current domain guide status enum — consider adding it or using a flag on the `diagnosing` status.
- **Cancellation** can only be performed by **Admin** per domain rules (not directly by customer).
- **Notifications** are referenced but the delivery mechanism (email, SMS, in-app) is not yet defined.
- **Deferred auth (Story 7a)** — booking form is accessible to guests. Form data is persisted to `localStorage` under a known key (e.g. `pending_booking`). After auth callback, the frontend checks for this key and auto-populates the form. Social login (Google) is presented as the primary/recommended option; email+password is secondary. The booking API endpoint still requires authentication — the form simply collects data client-side until auth is established.