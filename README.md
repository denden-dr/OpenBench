# OpenBench (PhoneFix Admin)

A clean, high-performance web application designed as a single-user administrative tool for a phone repair business. It simplifies device tracking, mirroring your Google Sheet workflow directly.

## Project Overview

OpenBench Admin provides a unified workspace for the shop owner to manage the entire intake, repair, payment, and warranty lifecycle of customer devices in one dashboard.

### Core Features

*   **Intake Form**: Quickly log repairs with customer details, device specifications (Brand/Model), accessories left (SIM, Case, etc.), price estimation, and warranty duration.
*   **Simple Status Workflow**: Track jobs through four linear stages:
    1.  **In for Service** (`service_in`): Newly dropped-off devices.
    2.  **On Process** (`on_process`): Devices actively being repaired.
    3.  **Fixed** (`fixed`): Repair complete, waiting for customer collection.
    4.  **Picked Up** (`picked_up`): Device retrieved. This logs the exit date, marks the payment as paid, and calculates the warranty expiry date.
*   **KPI Summary Cards**: Live counts of tickets at each stage and today's total revenue.
*   **Instant Search & Filters**: Search tickets by customer name, brand, model, or issue.

## Tech Stack

*   **Frontend**: Built with **Svelte 5** (Runes mode) and SvelteKit for rapid, reactive UI rendering. Tailwind CSS provides clean, premium aesthetics.
*   **Backend API**: High-performance RESTful API powered by **Go Fiber**.
*   **Database**: **PostgreSQL** (via `sqlx`) to persist ticket entries.

## Getting Started

### Prerequisites
- Docker (for database/Supabase containers)
- Go (v1.22+)
- Node.js (v20+)

### Running Locally

1. **Start the database**:
   ```bash
   make db-up
   ```
2. **Apply migrations**:
   ```bash
   make migrate-up
   ```
3. **Run the Go backend**:
   ```bash
   make run-backend
   ```
4. **Run SvelteKit frontend**:
   ```bash
   cd apps/frontend
   npm run dev
   ```
   Open `http://localhost:5173` to access the admin dashboard.
