# Frontend Warranty Claims Implementation Plan (Decoupled Queue)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Modify the frontend at `/warranty` and the mock layers to support a decoupled queue system:
1. **Pendaftaran (Intake)**: Searches for tickets, inputs the new issue, and registers the claim under `waiting_inspection` (sends a POST request without inspection decision).
2. **Antrean (Queue)**: Renders a list of all claims with status `waiting_inspection`.
3. **Pemeriksaan (Inspection)**: Supports clicking "Setujui" or "Tolak" (Void) directly on a queue item to resolve the claim.

**Tech Stack:** Svelte 5 (Runes), Tailwind CSS v4, Lucide Svelte.

---

## Files to Create/Modify
- `apps/frontend/src/lib/mocks/mockData.ts` (Modify: Add active and expired warranty test data, plus a pending claims seed array)
- `apps/frontend/src/lib/mocks/handlers.ts` (Modify: Implement GET, POST, approve, and void mock endpoints for claims)
- `apps/frontend/src/routes/warranty/+page.svelte` (Modify: Rebuild layout to show Search/Intake form at the top, and the Queue Table at the bottom)
- `apps/frontend/src/routes/+page.svelte` (Modify: Add navigation to warranty page, display warranty claims tags and active warranty status indicators)

---

### Task 1: Re-design Test Data & Mock Endpoints

**Files:**
- Modify: `apps/frontend/src/lib/mocks/mockData.ts`
- Modify: `apps/frontend/src/lib/mocks/handlers.ts`

- [ ] **Step 1: Modify mockData.ts to add picked-up tickets and initial claim queue**

Ensure mock tickets has:
- `TCK-004` (Apple iPhone 15 Pro, exit_date: 9 days ago, warranty_days: 30) - Active
- `TCK-005` (Xiaomi Mi 11 Ultra, exit_date: 39 days ago, warranty_days: 15) - Expired

Add a local mutable array for `mockWarrantyClaims`:
```typescript
export interface MockClaim {
  id: string;
  ticket_id: string;
  claim_ticket_id: string | null;
  issue: string;
  additional_description: string;
  status: 'waiting_inspection' | 'approved' | 'void';
  void_reason: string | null;
  inspected_at: string | null;
  created_at: string;
}

export let mockWarrantyClaims: MockClaim[] = [
  // Seed with one pending claim if desired, or leave empty
];
```

- [ ] **Step 2: Update mock handlers to support the decoupled queue endpoints**

In `apps/frontend/src/lib/mocks/handlers.ts`, implement:
- `POST /api/v1/warranty-claims` (Intake): Reads `ticket_id`, `issue`, `additional_description`. Inserts a claim with status `waiting_inspection`.
- `GET /api/v1/warranty-claims`: Returns list of claims.
- `POST /api/v1/warranty-claims/:id/approve`: Spawns a new ticket in `on_process` with price 0, `is_warranty = true`, and links it.
- `POST /api/v1/warranty-claims/:id/void`: Sets status to `void`, saves `void_reason`.

- [ ] **Step 3: Verify frontend type check compiles successfully**

Run: `cd apps/frontend && npm run check`

---

### Task 2: Build the Decoupled `/warranty` Page UI

**Files:**
- Modify: `apps/frontend/src/routes/warranty/+page.svelte`

- [ ] **Step 1: Re-write Svelte 5 logic & styling**

Rebuild the `/warranty/+page.svelte` file. It must show:
1. **Intake Section**:
   - Search bar.
   - Result card showing verification (Valid / Expired).
   - If Valid, a text input for **Kerusakan / Keluhan Baru** and a textarea for **Catatan**.
   - Button: **"Daftarkan Antrean Inspeksi"**. Clicking this calls `POST /api/v1/warranty-claims` and inserts it into the queue.
2. **Queue Section**:
   - Displays a clean table/list: **"Daftar Antrean Inspeksi Klaim"**.
   - Lists pending claims (Brand, Model, Customer, Issue, Created At).
   - Action columns: **"Setujui (Approve)"** and **"Tolak (Void)"** buttons.
   - Setujui calls the `/approve` endpoint.
   - Tolak opens a small dialog or prompt for `void_reason`, then calls the `/void` endpoint.

- [ ] **Step 2: Verify compiling the new Svelte page**

Run: `cd apps/frontend && npm run check`

---

### Task 3: Verification Checklist

Verify the decoupled flow:
- [ ] **Step 1: Pendaftaran Klaim (Intake)**
  - Open `/warranty`. Search for `TCK-004`.
  - Type issue: `"Layar Berkedip"`. Click **"Daftarkan Antrean Inspeksi"**.
  - Verify that the card resets, and the item now appears in the **"Daftar Antrean Inspeksi"** table below.
- [ ] **Step 2: Setujui Klaim (Approve)**
  - Locate the claim in the queue. Click **"Setujui"**.
  - Verify the claim disappears from the queue.
  - Return to the Dashboard and verify the new `TCK-W...` ticket is added with status `on_process` ("Klaim Diproses"), `Rp 0` price, and the purple warranty claim badge.
- [ ] **Step 3: Tolak Klaim (Void)**
  - Create another claim for `TCK-004` (issue: `"Baterai Kembung"`).
  - Locate it in the queue. Click **"Tolak"**.
  - Enter void reason: `"Ditemukan indikasi kemasukan air"`. Click submit.
  - Verify the claim is removed from the queue.
  - Return to the Dashboard and verify the new `TCK-W...` ticket appears with status `cancelled` ("Klaim Void"), `Rp 0` price, and the void reason in the notes.
