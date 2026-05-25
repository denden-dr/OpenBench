# Graph Report - OpenBench  (2026-05-25)

## Corpus Check
- 59 files · ~94,677 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 542 nodes · 634 edges · 46 communities (36 shown, 10 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 75 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9246912c`
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

## God Nodes (most connected - your core abstractions)
1. `TicketIntegrationTestSuite` - 27 edges
2. `str` - 16 edges
3. `Ticket` - 12 edges
4. `DesignSystemGenerator` - 11 edges
5. `compilerOptions` - 11 edges
6. `Repository Guidelines` - 10 edges
7. `_search_csv()` - 9 edges
8. `generate_design_system()` - 9 edges
9. `main()` - 9 edges
10. `NewPostgresStorage()` - 9 edges

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

## Communities (46 total, 10 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.09
Nodes (29): bool, DesignSystemGenerator, _detect_page_type(), format_ascii_box(), format_markdown(), format_master_md(), format_page_override_md(), generate_design_system() (+21 more)

### Community 1 - "Community 1"
Cohesion: 0.06
Nodes (29): dependencies, lucide-svelte, @supabase/supabase-js, devDependencies, svelte, svelte-check, @sveltejs/adapter-node, @sveltejs/kit (+21 more)

### Community 2 - "Community 2"
Cohesion: 0.08
Nodes (7): MockTicketRepository, MockTicketRepository_Create_Call, MockTicketRepository_Delete_Call, MockTicketRepository_Expecter, MockTicketRepository_GetByID_Call, MockTicketRepository_List_Call, MockTicketRepository_Update_Call

### Community 3 - "Community 3"
Cohesion: 0.06
Nodes (5): NewPostgresStorage(), WithCleanupInterval(), PostgresStorage, StorageOption, TicketIntegrationTestSuite

### Community 4 - "Community 4"
Cohesion: 0.09
Nodes (27): Backend Layered Architecture, Conventional Commits, Expand/Contract Migration Strategy, Graphify Knowledge Graph, Repository Guidelines, Mockery Mocks, Single-Instance Deployment, Testcontainers/PostgreSQL (+19 more)

### Community 5 - "Community 5"
Cohesion: 0.10
Nodes (19): main(), NewTicketHandler(), TicketHandler, hashIdempotencyRequest(), NewTicketIdempotency(), ScopeTicketIdempotencyKey(), TestHashIdempotencyRequest(), TestHashIdempotencyRequestChangesAfterValidationCorrection() (+11 more)

### Community 6 - "Community 6"
Cohesion: 0.13
Nodes (9): TestTicketErrorsAreStable(), TestValidateTicketUpdate(), ValidateTicketUpdate(), AppError, MapModelError(), MapRepositoryError(), TestMapModelError(), TestMapRepositoryError() (+1 more)

### Community 7 - "Community 7"
Cohesion: 0.15
Nodes (15): BM25, detect_domain(), _load_csv(), Lowercase, split, remove punctuation, filter short words, Build BM25 index from documents, Score all documents against query, Load CSV and return list of dicts, Core search function using BM25 (+7 more)

### Community 8 - "Community 8"
Cohesion: 0.16
Nodes (15): Config, DefaultDatabaseConfig(), getEnv(), getEnvDuration(), getEnvInt(), Load(), loadDatabaseConfig(), DatabaseConfig (+7 more)

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
Cohesion: 0.23
Nodes (4): MapDatabaseError(), TestMapDatabaseError(), sqlTicketRepository, TicketRepository

### Community 13 - "Community 13"
Cohesion: 0.29
Nodes (7): Electronics Repair Workstation, Fine Tweezers, Hot Air Rework Station, Inspection Microscope, Logic Board, Precision Screwdriver Set, Repair Mat

### Community 14 - "Community 14"
Cohesion: 0.06
Nodes (31): Accessibility, Available Domains, Available Stacks, code:bash (python3 --version || python --version), code:bash (python3 skills/ui-ux-pro-max/scripts/search.py "beauty spa w), code:bash (# Get UX guidelines for animation and accessibility), code:bash (python3 skills/ui-ux-pro-max/scripts/search.py "layout respo), code:bash (# ASCII box (default) - best for terminal display) (+23 more)

### Community 15 - "Community 15"
Cohesion: 0.48
Nodes (4): handleMockRequest(), mockTickets, setMockTickets(), handle()

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

### Community 40 - "Community 40"
Cohesion: 0.17
Nodes (12): code:bash (python3 skills/ui-ux-pro-max/scripts/search.py "<keyword>" -), code:bash (python3 skills/ui-ux-pro-max/scripts/search.py "<product_typ), code:bash (python3 skills/ui-ux-pro-max/scripts/search.py "beauty spa w), code:bash (python3 skills/ui-ux-pro-max/scripts/search.py "<query>" --d), code:bash (python3 skills/ui-ux-pro-max/scripts/search.py "<query>" --d), code:bash (python3 skills/ui-ux-pro-max/scripts/search.py "<keyword>" -), How to Use This Skill, Step 1: Analyze User Requirements (+4 more)

### Community 41 - "Community 41"
Cohesion: 0.18
Nodes (10): Build, Test, and Development Commands, Coding Style & Naming Conventions, Commit & Pull Request Guidelines, Deployment & Schema Migrations, graphify, Project Structure & Module Organization, Repository Guidelines, Runner Notes (+2 more)

### Community 42 - "Community 42"
Cohesion: 0.20
Nodes (9): 1. Check for `.gitignore`, 2. Analyze the Context, 3. Stage Files Selectively, 4. Commit with a Standardized Message, code:block1 (<type>[optional scope]: <description>), 🚫 CRITICAL RULE: NEVER USE `git add .`, Example, Git Commit Skill (+1 more)

### Community 43 - "Community 43"
Cohesion: 0.14
Nodes (13): Assessment, Code Review: PR 17 Fixes, Code Review: PR 17 - Idempotency and Loading-Aware Ticket Flows, Critical, Critical (Must Fix), DONE, Important, Important (Should Fix) (+5 more)

### Community 44 - "Community 44"
Cohesion: 0.22
Nodes (8): Building, code:sh (# create a new project), code:sh (# recreate this project), code:sh (npm run dev), code:sh (npm run build), Creating a project, Developing, sv

### Community 45 - "Community 45"
Cohesion: 0.29
Nodes (6): code:bash (rtk git status), code:bash (rtk gain              # Show token savings), Meta Commands, RTK - Rust Token Killer (Google Antigravity), Rule, Why

## Knowledge Gaps
- **203 isolated node(s):** `bool`, `extends`, `rewriteRelativeImportExtensions`, `allowJs`, `checkJs` (+198 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **10 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `ScopeTicketIdempotencyKey()` connect `Community 5` to `Community 6`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Why does `TicketIntegrationTestSuite` connect `Community 3` to `Community 8`, `Community 5`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Why does `main()` connect `Community 5` to `Community 8`, `Community 3`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Are the 3 inferred relationships involving `str` (e.g. with `.tokenize()` and `_search_csv()`) actually correct?**
  _`str` has 3 INFERRED edges - model-reasoned connections that need verification._
- **What connects `BM25 ranking algorithm for text search`, `Lowercase, split, remove punctuation, filter short words`, `Build BM25 index from documents` to the rest of the system?**
  _229 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 0` be split into smaller, more focused modules?**
  _Cohesion score 0.08771929824561403 - nodes in this community are weakly interconnected._
- **Should `Community 1` be split into smaller, more focused modules?**
  _Cohesion score 0.058823529411764705 - nodes in this community are weakly interconnected._