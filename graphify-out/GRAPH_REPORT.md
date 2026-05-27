# Graph Report - OpenBench  (2026-05-27)

## Corpus Check
- 107 files · ~117,858 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1118 nodes · 1459 edges · 77 communities (53 shown, 24 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 161 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c1f9f0b3`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 17|Community 17]]
- [[_COMMUNITY_Community 18|Community 18]]
- [[_COMMUNITY_Community 19|Community 19]]
- [[_COMMUNITY_Community 20|Community 20]]
- [[_COMMUNITY_Community 21|Community 21]]
- [[_COMMUNITY_Community 22|Community 22]]
- [[_COMMUNITY_Community 23|Community 23]]
- [[_COMMUNITY_Community 29|Community 29]]
- [[_COMMUNITY_Community 30|Community 30]]
- [[_COMMUNITY_Community 31|Community 31]]
- [[_COMMUNITY_Community 32|Community 32]]
- [[_COMMUNITY_Community 33|Community 33]]
- [[_COMMUNITY_Community 34|Community 34]]
- [[_COMMUNITY_Community 35|Community 35]]
- [[_COMMUNITY_Community 36|Community 36]]
- [[_COMMUNITY_Community 37|Community 37]]
- [[_COMMUNITY_Community 38|Community 38]]
- [[_COMMUNITY_Community 39|Community 39]]
- [[_COMMUNITY_Community 40|Community 40]]
- [[_COMMUNITY_Community 41|Community 41]]
- [[_COMMUNITY_Community 42|Community 42]]
- [[_COMMUNITY_Community 43|Community 43]]
- [[_COMMUNITY_Community 44|Community 44]]
- [[_COMMUNITY_Community 45|Community 45]]
- [[_COMMUNITY_Community 46|Community 46]]
- [[_COMMUNITY_Community 47|Community 47]]
- [[_COMMUNITY_Community 48|Community 48]]
- [[_COMMUNITY_Community 49|Community 49]]
- [[_COMMUNITY_Community 50|Community 50]]
- [[_COMMUNITY_Community 51|Community 51]]
- [[_COMMUNITY_Community 52|Community 52]]
- [[_COMMUNITY_Community 53|Community 53]]
- [[_COMMUNITY_Community 54|Community 54]]
- [[_COMMUNITY_Community 55|Community 55]]
- [[_COMMUNITY_Community 56|Community 56]]
- [[_COMMUNITY_Community 57|Community 57]]
- [[_COMMUNITY_Community 58|Community 58]]
- [[_COMMUNITY_Community 59|Community 59]]
- [[_COMMUNITY_Community 60|Community 60]]
- [[_COMMUNITY_Community 61|Community 61]]
- [[_COMMUNITY_Community 62|Community 62]]
- [[_COMMUNITY_Community 63|Community 63]]
- [[_COMMUNITY_Community 64|Community 64]]
- [[_COMMUNITY_Community 65|Community 65]]
- [[_COMMUNITY_Community 70|Community 70]]
- [[_COMMUNITY_Community 71|Community 71]]
- [[_COMMUNITY_Community 72|Community 72]]
- [[_COMMUNITY_Community 76|Community 76]]
- [[_COMMUNITY_Community 77|Community 77]]
- [[_COMMUNITY_Community 78|Community 78]]
- [[_COMMUNITY_Community 79|Community 79]]
- [[_COMMUNITY_Community 80|Community 80]]
- [[_COMMUNITY_Community 81|Community 81]]

## God Nodes (most connected - your core abstractions)
1. `TicketIntegrationTestSuite` - 28 edges
2. `WarrantyClaimIntegrationTestSuite` - 24 edges
3. `MapDatabaseError()` - 20 edges
4. `main()` - 18 edges
5. `str` - 16 edges
6. `Infrastructure` - 16 edges
7. `MapRepositoryError()` - 14 edges
8. `MockTicketRepository` - 13 edges
9. `Ticket` - 12 edges
10. `sqlTicketRepository` - 12 edges

## Surprising Connections (you probably didn't know these)
- `Repository Guidelines` --references--> `RTK (Rust Token Killer)`  [EXTRACTED]
  AGENTS.md → .agents/rules/antigravity-rtk-rules.md
- `Repair Status Workflow` --references--> `Status: service_in`  [EXTRACTED]
  README.md → docs/PRD.md
- `Repair Status Workflow` --references--> `Status: on_process`  [EXTRACTED]
  README.md → docs/PRD.md
- `Repair Status Workflow` --references--> `Status: fixed`  [EXTRACTED]
  README.md → docs/PRD.md
- `Repair Status Workflow` --references--> `Status: picked_up`  [EXTRACTED]
  README.md → docs/PRD.md

## Hyperedges (group relationships)
- **Phone Repair Bench Setup** — equipment_inspection_microscope, equipment_hot_air_rework_station, equipment_precision_screwdriver_set, equipment_fine_tweezers [INFERRED 0.93]
- **Phone Board Microsoldering Workflow** — microsoldering_repair_scene, smartphone_logic_board, soldering_iron, precision_tweezers, a15_bionic_chip [INFERRED 0.92]
- **Repair Status Workflow Stages** — prd_ticket_status_service_in, prd_ticket_status_on_process, prd_ticket_status_fixed, prd_ticket_status_picked_up [EXTRACTED 1.00]

## Communities (77 total, 24 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.06
Nodes (44): bool, BM25, detect_domain(), _load_csv(), Lowercase, split, remove punctuation, filter short words, Build BM25 index from documents, Score all documents against query, Load CSV and return list of dicts (+36 more)

### Community 1 - "Community 1"
Cohesion: 0.24
Nodes (10): stubResult, NewTicketService(), ptrInt(), TestTicketService_CreateTicket(), TestTicketService_DeleteTicket(), TestTicketService_GetPublicTicket(), TestTicketService_GetTicket(), TestTicketService_ListTickets() (+2 more)

### Community 2 - "Community 2"
Cohesion: 0.06
Nodes (12): MockTicketRepository, MockTicketRepository_BeginTx_Call, MockTicketRepository_Create_Call, MockTicketRepository_CreateTx_Call, MockTicketRepository_Delete_Call, MockTicketRepository_GetByID_Call, MockTicketRepository_GetByIDForUpdateTx_Call, MockTicketRepository_GetByIDs_Call (+4 more)

### Community 3 - "Community 3"
Cohesion: 0.06
Nodes (5): NewPostgresStorage(), WithCleanupInterval(), PostgresStorage, StorageOption, TicketIntegrationTestSuite

### Community 4 - "Community 4"
Cohesion: 0.09
Nodes (27): Backend Layered Architecture, Conventional Commits, Expand/Contract Migration Strategy, Graphify Knowledge Graph, Repository Guidelines, Mockery Mocks, Single-Instance Deployment, Testcontainers/PostgreSQL (+19 more)

### Community 5 - "Community 5"
Cohesion: 0.31
Nodes (8): SetupTestDB(), NewHealthHandler(), RegisterRoutes(), NewTicketHandler(), TestTicketHandler_GetByID(), NewWarrantyClaimHandler(), NewTicketRepository(), NewWarrantyClaimRepository()

### Community 6 - "Community 6"
Cohesion: 0.08
Nodes (19): TestTicketErrorsAreStable(), TestValidateTicketUpdate(), ValidateTicketUpdate(), AppError, MapModelError(), MapRepositoryError(), TestMapModelError(), TestMapRepositoryError() (+11 more)

### Community 7 - "Community 7"
Cohesion: 0.06
Nodes (9): MockWarrantyClaimRepository, MockWarrantyClaimRepository_BeginTx_Call, MockWarrantyClaimRepository_Create_Call, MockWarrantyClaimRepository_Expecter, MockWarrantyClaimRepository_GetByID_Call, MockWarrantyClaimRepository_GetByIDForUpdateTx_Call, MockWarrantyClaimRepository_GetOpenClaimByTicketID_Call, MockWarrantyClaimRepository_List_Call (+1 more)

### Community 8 - "Community 8"
Cohesion: 0.08
Nodes (19): Config, DefaultDatabaseConfig(), getEnv(), getEnvDuration(), getEnvInt(), Load(), loadDatabaseConfig(), TestLoad_DefaultCORSAllowOrigins() (+11 more)

### Community 9 - "Community 9"
Cohesion: 0.20
Nodes (5): Ticket, optionalText(), TicketPaymentStatus, TicketStatus, TicketUpdate

### Community 10 - "Community 10"
Cohesion: 0.14
Nodes (13): files, code, document, image, paper, video, graphifyignore_patterns, needs_graph (+5 more)

### Community 11 - "Community 11"
Cohesion: 0.15
Nodes (12): compilerOptions, allowJs, checkJs, esModuleInterop, forceConsistentCasingInFileNames, moduleResolution, resolveJsonModule, rewriteRelativeImportExtensions (+4 more)

### Community 12 - "Community 12"
Cohesion: 0.11
Nodes (6): MapDatabaseError(), TestMapDatabaseError(), sqlTicketRepository, sqlWarrantyClaimRepository, TicketRepository, WarrantyClaimRepository

### Community 13 - "Community 13"
Cohesion: 0.29
Nodes (7): Electronics Repair Workstation, Fine Tweezers, Hot Air Rework Station, Inspection Microscope, Logic Board, Precision Screwdriver Set, Repair Mat

### Community 14 - "Community 14"
Cohesion: 0.05
Nodes (43): Accessibility, Available Domains, Available Stacks, code:bash (python3 --version || python --version), code:bash (python3 skills/ui-ux-pro-max/scripts/search.py "<keyword>" -), code:bash (python3 skills/ui-ux-pro-max/scripts/search.py "beauty spa w), code:bash (# Get UX guidelines for animation and accessibility), code:bash (python3 skills/ui-ux-pro-max/scripts/search.py "layout respo) (+35 more)

### Community 15 - "Community 15"
Cohesion: 0.36
Nodes (7): handleMockRequest(), MockClaim, mockTickets, mockWarrantyClaims, setMockTickets(), setMockWarrantyClaims(), handle()

### Community 16 - "Community 16"
Cohesion: 0.33
Nodes (6): A15 Bionic Chip, Hero Image, Microsoldering Repair Scene, Precision Tweezers, Smartphone Logic Board, Soldering Iron

### Community 17 - "Community 17"
Cohesion: 0.33
Nodes (5): content, file, fs, path, replacements

### Community 18 - "Community 18"
Cohesion: 0.50
Nodes (3): CreateTicketRequest, TicketResponse, UpdateTicketRequest

### Community 21 - "Community 21"
Cohesion: 0.67
Nodes (3): OpenBench Favicon, Repair Tools, Smartphone Outline

### Community 35 - "Community 35"
Cohesion: 0.14
Nodes (13): 1. Project Overview, 2. User Roles & Permissions, 3.1 Repair Ticket Intake & Form, 3.2 Main Repair Admin Dashboard, 3.3 Status & Workflows, 3. Functional Requirements, 4.1 Frontend (Svelte), 4.2 Backend (Go + Fiber) (+5 more)

### Community 36 - "Community 36"
Cohesion: 0.14
Nodes (13): 1. Prasyarat Pengujian, 2. Skenario Pengujian, code:bash (make compose-up), code:bash (make migrate-up), code:bash (make run-backend), code:bash (make run-frontend), code:bash (curl -X POST http://localhost:3000/api/v1/tickets \), code:bash (curl -X POST http://localhost:3000/api/v1/tickets \) (+5 more)

### Community 37 - "Community 37"
Cohesion: 0.15
Nodes (12): 1. Create ticket, 2. Idempotency key regeneration, 3. Duplicate request handling, 4. Loading state, 5. Quick action, 6. Delete / cancel flow, Issues Found, Manual Test Report (+4 more)

### Community 38 - "Community 38"
Cohesion: 0.15
Nodes (12): code:block1 (apps/), code:bash (# 1. Start database + services), code:bash (graphify update .), Core Features, Getting Started, Makefile Targets, OpenBench (PhoneFix Admin), Prerequisites (+4 more)

### Community 39 - "Community 39"
Cohesion: 0.17
Nodes (11): Buttons and Actions, Common Mistakes & Red Flags, Core Principles, Empty States, Error Messages, Implementation Examples, Overview, Quick Reference (+3 more)

### Community 41 - "Community 41"
Cohesion: 0.18
Nodes (10): Build, Test, and Development Commands, Coding Style & Naming Conventions, Commit & Pull Request Guidelines, Deployment & Schema Migrations, graphify, Project Structure & Module Organization, Repository Guidelines, Runner Notes (+2 more)

### Community 42 - "Community 42"
Cohesion: 0.20
Nodes (9): 1. Check for `.gitignore`, 2. Analyze the Context, 3. Stage Files Selectively, 4. Commit with a Standardized Message, code:block1 (<type>[optional scope]: <description>), 🚫 CRITICAL RULE: NEVER USE `git add .`, Example, Git Commit Skill (+1 more)

### Community 43 - "Community 43"
Cohesion: 0.07
Nodes (40): 1. Edit form uses `PUT`, but backend only registers `PATCH`, 1. Race condition: pengecekan status klaim melewati isolasi transaksi, 1. Void/cancelled warranty tickets mendapat 30-day warranty, 2. Logika insert tiket duplikat di warranty claim repo, 2. Quick status action calls an unregistered `/status` route, 2. Validation test false positive, 3. Bentuk response tidak konsisten pada Approve/Void vs List/Create, 3. N+1 query di `ListClaims` untuk ticket enrichment (+32 more)

### Community 44 - "Community 44"
Cohesion: 0.22
Nodes (8): Building, code:sh (# create a new project), code:sh (# recreate this project), code:sh (npm run dev), code:sh (npm run build), Creating a project, Developing, sv

### Community 45 - "Community 45"
Cohesion: 0.29
Nodes (6): code:bash (rtk git status), code:bash (rtk gain              # Show token savings), Meta Commands, RTK - Rust Token Killer (Google Antigravity), Rule, Why

### Community 46 - "Community 46"
Cohesion: 0.06
Nodes (33): Backend Warranty Claims Implementation Plan, Backend Warranty Claims Implementation Plan (Decoupled Queue), code:sql (ALTER TABLE tickets ADD COLUMN is_warranty BOOLEAN NOT NULL ), code:go (package repository), code:bash (git add apps/backend/internal/model/warranty_claim.go apps/b), code:go (IsWarranty            *bool            `json:"is_warranty" v), code:go (IsWarranty            bool            `json:"is_warranty"`), code:go (package dto) (+25 more)

### Community 47 - "Community 47"
Cohesion: 0.09
Nodes (24): code:typescript (export interface MockClaim {), code:html ({getStatusLabel(ticket.status, Boolean(ticket.is_warranty))}), code:html (<a), code:html ({#if ticket.is_warranty}), code:bash (git add apps/frontend/src/routes/+page.svelte), code:typescript (// Matches /api/v1/warranty-claims), code:bash (git add apps/frontend/src/lib/mocks/mockData.ts apps/fronten), code:html (<script lang="ts">) (+16 more)

### Community 48 - "Community 48"
Cohesion: 0.12
Nodes (20): 1. Overview & Business Logic, 2. Technical Architecture, 3. Frontend UI Flow (Svelte 5), 4. Sequence Diagram, A. Database Schema, A. Dedicated `/warranty` Page Layout, B. Backend Go Models (`apps/backend/internal/model`), B. Backend Go Models (`apps/backend/internal/model/warranty_claim.go`) (+12 more)

### Community 49 - "Community 49"
Cohesion: 0.40
Nodes (4): ClaimCreationResult, CreateWarrantyClaimRequest, VoidWarrantyClaimRequest, WarrantyClaimResponse

### Community 51 - "Community 51"
Cohesion: 0.17
Nodes (11): 1. Arsitektur Alur Kerja Baru, 2. Perubahan Database & Migrasi, 3. Implementasi Kode Backend (Go Fiber), 4. Implementasi Frontend (Svelte 5 & SvelteKit), 5. Pengujian & Verifikasi, 6. Tindakan Lanjut Temuan Code Review (issue.md), A. Model (`apps/backend/internal/model/`), B. Repository (`apps/backend/internal/repository/`) (+3 more)

### Community 52 - "Community 52"
Cohesion: 0.25
Nodes (7): 1. Edit form uses `PUT`, but backend only registers `PATCH`, 2. Quick status action calls an unregistered `/status` route, 3. Technician issue and customer approval/rejection actions call missing routes, Important (Should Fix), Issues, Minor (Nice to Have), PR #20 Code Review Issues

### Community 55 - "Community 55"
Cohesion: 0.40
Nodes (4): Deskripsi Masalah, Detail Lokasi Kode, Rencana: Implementasi Pagination pada Endpoint List, Rencana Solusi

### Community 57 - "Community 57"
Cohesion: 0.06
Nodes (31): warranty, warranty, svelte, ../app.css, $lib/types/ticket, $lib/utils/format, $lib/utils/warranty, activeIndex (+23 more)

### Community 59 - "Community 59"
Cohesion: 0.06
Nodes (9): MockTicketService, MockTicketService_CreateTicket_Call, MockTicketService_DeleteTicket_Call, MockTicketService_Expecter, MockTicketService_GetPublicTicket_Call, MockTicketService_GetTicket_Call, MockTicketService_ListTickets_Call, MockTicketService_TrackPublicTicket_Call (+1 more)

### Community 60 - "Community 60"
Cohesion: 0.06
Nodes (34): Backend, Frontend, Infrastructure, L-10. Bad Practice — `handlers.ts` event typed as `any`, L-11. Bad Practice — Multiple a11y ignore comments, L-12. UX — Bahasa campuran di step descriptions track page, L-13. UX — Stats section horizontal scroll di mobile, L-14. Performance — `Intl.NumberFormat` di-construct setiap pemanggilan (+26 more)

### Community 61 - "Community 61"
Cohesion: 0.08
Nodes (7): MockTransaction, MockTransaction_Commit_Call, MockTransaction_ExecContext_Call, MockTransaction_Expecter, MockTransaction_GetContext_Call, MockTransaction_QueryRowxContext_Call, MockTransaction_Rollback_Call

### Community 62 - "Community 62"
Cohesion: 0.09
Nodes (6): MockWarrantyClaimService, MockWarrantyClaimService_ApproveClaim_Call, MockWarrantyClaimService_CreateClaim_Call, MockWarrantyClaimService_Expecter, MockWarrantyClaimService_ListClaims_Call, MockWarrantyClaimService_VoidClaim_Call

### Community 65 - "Community 65"
Cohesion: 0.13
Nodes (15): 1. Status Perbaikan Sebelumnya, 1. Status Perbaikan Sebelumnya (Resolved), 1. Status Perbaikan Secara Keseluruhan (All Resolved 🎉), 1. Unhandled Error pada `tx.Rollback()` (Bug / Best Practice), 2. Status Pengembangan Lanjutan, 2. Temuan Isu Baru, 2. Tidak Ada Graceful Shutdown pada Fiber App (Best Practice / Konsistensi), 3. Konfigurasi CORS Belum Didefinisikan (Potensi Bug Integrasi) (+7 more)

### Community 76 - "Community 76"
Cohesion: 0.21
Nodes (14): main(), hashIdempotencyRequest(), idempotencyConcretePath(), NewIdempotency(), NewTicketIdempotency(), ScopeIdempotencyKey(), ScopeTicketIdempotencyKey(), TestHashIdempotencyRequest() (+6 more)

### Community 77 - "Community 77"
Cohesion: 0.07
Nodes (26): dependencies, lucide-svelte, qrcode, @supabase/supabase-js, devDependencies, svelte-check, @sveltejs/adapter-node, @sveltejs/kit (+18 more)

### Community 79 - "Community 79"
Cohesion: 0.53
Nodes (5): NewWarrantyClaimService(), TestWarrantyClaimService_ApproveClaim(), TestWarrantyClaimService_CreateClaim(), TestWarrantyClaimService_ListClaims(), TestWarrantyClaimService_VoidClaim()

## Knowledge Gaps
- **355 isolated node(s):** `bool`, `extends`, `rewriteRelativeImportExtensions`, `allowJs`, `checkJs` (+350 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **24 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `main()` connect `Community 76` to `Community 1`, `Community 3`, `Community 5`, `Community 8`, `Community 79`, `Community 81`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Why does `ScopeIdempotencyKey()` connect `Community 76` to `Community 5`, `Community 6`?**
  _High betweenness centrality (0.010) - this node is a cross-community bridge._
- **Why does `TicketIntegrationTestSuite` connect `Community 3` to `Community 8`, `Community 5`?**
  _High betweenness centrality (0.010) - this node is a cross-community bridge._
- **Are the 19 inferred relationships involving `MapDatabaseError()` (e.g. with `TestMapDatabaseError()` and `.Is()`) actually correct?**
  _`MapDatabaseError()` has 19 INFERRED edges - model-reasoned connections that need verification._
- **Are the 17 inferred relationships involving `main()` (e.g. with `Load()` and `NewDB()`) actually correct?**
  _`main()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 3 inferred relationships involving `str` (e.g. with `.tokenize()` and `_search_csv()`) actually correct?**
  _`str` has 3 INFERRED edges - model-reasoned connections that need verification._
- **What connects `BM25 ranking algorithm for text search`, `Lowercase, split, remove punctuation, filter short words`, `Build BM25 index from documents` to the rest of the system?**
  _381 weakly-connected nodes found - possible documentation gaps or missing edges._