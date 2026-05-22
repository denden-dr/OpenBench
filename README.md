# OpenBench (PhoneFix)

A comprehensive, end-to-end web application designed to manage the workflow of a phone repair business. The system serves both customers (booking and tracking) and staff (workflow and inventory management).

## Project Overview

OpenBench is built to handle the entire lifecycle of device repairs, tracking everything from the initial customer drop-off to technician assignment, parts usage, invoicing, and final payment. 

### Key Features by Role

*   **Guest / Public**: Ability to track repair status via a Ticket ID and view service pricing.
*   **Customer**: Can create a profile, book new repairs via a multi-step intake form (selecting device, brand, model, and issue), and view repair history. The intake process captures device access details and includes a mandatory diagnosis fee agreement.
*   **Technician**: Can claim tickets from a queue, diagnose issues, document repairs with "before" and "after" media, log parts used from inventory, and manage the repair status flow.
*   **Admin**: Full access to manage parts inventory (tracking ODM vs Original grades), view financial reporting, and oversee all system operations.

## Tech Stack

OpenBench utilizes a modern, hybrid architecture:

*   **Frontend**: Built with **SvelteKit**, providing a hybrid SSR/SPA experience. Styling relies on Vanilla CSS and Modern CSS Variables to ensure a high-performance, premium UI.
*   **Backend API**: Powered by **Go Fiber**, providing a robust, high-performance RESTful API.
*   **Database**: **PostgreSQL** is the source of truth for business data, managing tickets, inventory, users, and payments.
*   **Authentication**: **Supabase Auth** is used for secure identity management. The system employs a secure, server-side HTTP-only cookie approach (`@supabase/ssr` on SvelteKit) to mitigate XSS risks, with the Go backend validating Supabase JWTs.

## Architecture Highlights

1.  **Secure Authentication Flow**: 
    The client logs in via SvelteKit, which communicates with Supabase. Supabase returns a session, and SvelteKit sets a secure `HttpOnly` cookie containing the JWT. All subsequent API calls to the Go Fiber backend include this cookie, which is extracted and verified before executing business logic.
2.  **Concurrency Management**: 
    Ticket claiming operates on a first-come, first-served basis using atomic database updates to prevent race conditions when multiple technicians attempt to claim the same repair ticket.
3.  **Inventory Transactions**:
    Part deduction and status updates are wrapped in ACID database transactions, rolling back if stock is insufficient.
4.  **Third-Party Integrations**:
    Supports asynchronous webhook integrations (like Midtrans) to handle online payments and automatically update ticket/payment statuses upon success.

## Documentation

For deeper dives into the project's requirements and technical design, please refer to the documentation in the `docs/` directory:

*   **[PRD.md](./docs/PRD.md)**: Product Requirements Document, including detailed roles, feature specifications, and the database schema.
*   **[ARCHITECTURE.md](./docs/ARCHITECTURE.md)**: Detailed system architecture, sequence diagrams for core workflows (auth, ticket claiming, repair workflow, webhooks), and component responsibilities.
