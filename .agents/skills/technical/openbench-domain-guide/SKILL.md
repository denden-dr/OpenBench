---
name: openbench-domain-guide
description: Use when implementing any feature, writing queries, designing UI, or making decisions that involve repair tickets, user roles, inventory, payments, or business rules in the OpenBench (PhoneFix) project
---

# OpenBench Domain Guide

## Overview

OpenBench (PhoneFix) is a phone repair management system. This is the single source of truth for domain concepts, business rules, and data relationships. **Always consult this before writing any feature code.**

## Roles & Permissions

| Action | Guest | Customer | Technician | Admin |
|--------|:-----:|:--------:|:----------:|:-----:|
| Track ticket (by ID + phone) | ✅ | ✅ | ✅ | ✅ |
| View service prices | ✅ | ✅ | ✅ | ✅ |
| Book repair / create ticket | ❌ | ✅ | ❌ | ✅ |
| View own repair history | ❌ | ✅ | ❌ | ✅ |
| Download invoice/receipt | ❌ | ✅ | ❌ | ✅ |
| View unassigned ticket queue | ❌ | ❌ | ✅ | ✅ |
| Claim/take a ticket | ❌ | ❌ | ✅ | ✅ |
| Update ticket status (own) | ❌ | ❌ | ✅ | ✅ |
| Add diagnosis notes | ❌ | ❌ | ✅ | ✅ |
| Upload before/after photos | ❌ | ❌ | ✅ | ✅ |
| Add internal comments | ❌ | ❌ | ✅ | ✅ |
| Log parts used on ticket | ❌ | ❌ | ✅ | ✅ |
| Manage inventory/parts | ❌ | ❌ | ❌ | ✅ |
| Manage users | ❌ | ❌ | ❌ | ✅ |
| View reports/analytics | ❌ | ❌ | ❌ | ✅ |
| System settings | ❌ | ❌ | ❌ | ✅ |

## Ticket Lifecycle

```
received → diagnosing → waiting_parts → repairing → ready → completed
                    ↘                                         ↗
                      ─────────── cancelled ──────────────────
```

### State Transitions

| From | To | Triggered By | Conditions |
|------|----|-------------|------------|
| `received` | `diagnosing` | Technician claims ticket | `technician_id IS NULL` (first-come-first-served) |
| `diagnosing` | `waiting_parts` | Technician | Parts needed but not in stock |
| `diagnosing` | `repairing` | Technician | Parts available, logs parts used (inventory deducted in transaction) |
| `diagnosing` | `cancelled` | Admin | Customer declines after diagnosis |
| `waiting_parts` | `repairing` | Technician | Parts arrived, logs parts used |
| `repairing` | `ready` | Technician | Repair complete, "after" photos uploaded |
| `ready` | `completed` | System/Admin | Payment received (cash confirmed or webhook received) |
| Any (except `completed`) | `cancelled` | Admin only | With reason logged in audit |

### State Rules
- **Only forward transitions** — no going back (except cancel)
- **Claiming is atomic** — `UPDATE ... WHERE technician_id IS NULL` prevents double-claim (409 Conflict if already taken)
- **`estimated_ready_at`** — must be set during `diagnosing` phase, visible to customer
- **`warranty_expiry`** — calculated on `completed`: ODM parts = +7 days, Original parts = +30 days

## Business Rules

### Diagnosis Fee
- **Mandatory** — displayed before booking, customer must agree via checkbox
- Charged regardless of whether repair proceeds
- Included in final invoice: `Total = diagnosis_fee + labor_fee + Σ(parts)`

### Warranty Calculation
| Part Grade | Warranty Duration |
|-----------|-------------------|
| ODM | 7 days from completion |
| Original | 30 days from completion |
- If mixed parts used: warranty = **shortest** (7 days)
- `warranty_expiry` = `completed_at + min(warranty_durations)`

### Inventory
- Stock deducted **within a database transaction** alongside ticket update
- If stock insufficient → transaction rollback → 400 error
- Part grades: `ODM` (aftermarket) or `Original` (OEM)
- Audit log entry for every stock change

### Payments
- Methods: `cash` or `online` (Midtrans)
- Online payments are **async** — status updated via webhook
- Payment statuses: `pending → completed | failed`
- Ticket moves to `completed` only after payment `completed`

### Accessory Tracking
- Checklist items (SIM, Case, SD Card, etc.) stored as JSON on ticket
- Reminder shown to technician on status change to `ready`

## Data Model Quick Reference

### Primary Entities
| Entity | PK | Key Fields |
|--------|-----|-----------|
| `users` | UUID | `supabase_uid` (unique), `email`, `role` (admin/technician/customer), `name` |
| `customers` | UUID | `user_id` → users, `phone`, `address` |
| `tickets` | UUID/ShortID | `customer_id`, `device_type`, `brand`, `model`, `status`, `technician_id`, `diagnosis_fee`, `labor_fee` |
| `parts` | UUID | `name`, `grade` (ODM/Original), `price`, `stock_level`, `brand_compatibility` |
| `payments` | UUID | `ticket_id`, `amount`, `method`, `status`, `external_id` (Midtrans order ID) |

### Junction & Support Tables
| Entity | Purpose |
|--------|---------|
| `ticket_parts` | Parts used on a ticket (`ticket_id`, `part_id`, `quantity`) |
| `attachments` | Before/after photos & receipts (`ticket_id`, `file_url`, `type`) |
| `ticket_comments` | Internal tech notes (`ticket_id`, `user_id`, `content`) — **not visible to customer** |
| `audit_logs` | All sensitive actions (`user_id`, `action`, `entity_type`, `entity_id`, `changes` JSON) |

### Key Relationships
```
users 1──* customers (via user_id)
customers 1──* tickets (via customer_id)
users 1──* tickets (via technician_id, nullable)
tickets 1──* ticket_parts *──1 parts
tickets 1──* attachments
tickets 1──* ticket_comments
tickets 1──* payments
users 1──* audit_logs
```

## Domain Glossary

| Term | Meaning |
|------|---------|
| **Ticket** | A repair job from intake to completion |
| **Claim** | Technician self-assigns an unassigned ticket |
| **Diagnosis Fee** | Mandatory upfront fee for inspection |
| **Part Grade** | ODM (aftermarket) vs Original (OEM) |
| **ShortID** | Human-readable ticket identifier for customer tracking |
| **Before/After Photos** | Liability documentation uploaded by technician |
| **Internal Comment** | Tech-only notes, hidden from customer |
| **Repair Mode** | Device access method alternative to passcode |
