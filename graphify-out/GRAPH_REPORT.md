# Graph Report - OpenBench  (2026-06-22)

## Corpus Check
- 143 files · ~65,538 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 924 nodes · 1099 edges · 94 communities (61 shown, 33 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 105 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e9051715`
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
- [[_COMMUNITY_Community 24|Community 24]]
- [[_COMMUNITY_Community 25|Community 25]]
- [[_COMMUNITY_Community 26|Community 26]]
- [[_COMMUNITY_Community 27|Community 27]]
- [[_COMMUNITY_Community 28|Community 28]]
- [[_COMMUNITY_Community 29|Community 29]]
- [[_COMMUNITY_Community 30|Community 30]]
- [[_COMMUNITY_Community 31|Community 31]]
- [[_COMMUNITY_Community 32|Community 32]]
- [[_COMMUNITY_Community 33|Community 33]]
- [[_COMMUNITY_Community 34|Community 34]]
- [[_COMMUNITY_Community 35|Community 35]]
- [[_COMMUNITY_Community 36|Community 36]]
- [[_COMMUNITY_Community 37|Community 37]]
- [[_COMMUNITY_Community 39|Community 39]]
- [[_COMMUNITY_Community 40|Community 40]]
- [[_COMMUNITY_Community 41|Community 41]]
- [[_COMMUNITY_Community 42|Community 42]]
- [[_COMMUNITY_Community 43|Community 43]]
- [[_COMMUNITY_Community 44|Community 44]]
- [[_COMMUNITY_Community 45|Community 45]]
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
- [[_COMMUNITY_Community 65|Community 65]]
- [[_COMMUNITY_Community 66|Community 66]]
- [[_COMMUNITY_Community 67|Community 67]]
- [[_COMMUNITY_Community 68|Community 68]]
- [[_COMMUNITY_Community 69|Community 69]]
- [[_COMMUNITY_Community 70|Community 70]]
- [[_COMMUNITY_Community 71|Community 71]]
- [[_COMMUNITY_Community 72|Community 72]]
- [[_COMMUNITY_Community 73|Community 73]]
- [[_COMMUNITY_Community 74|Community 74]]
- [[_COMMUNITY_Community 75|Community 75]]
- [[_COMMUNITY_Community 76|Community 76]]
- [[_COMMUNITY_Community 77|Community 77]]
- [[_COMMUNITY_Community 78|Community 78]]
- [[_COMMUNITY_Community 79|Community 79]]
- [[_COMMUNITY_Community 80|Community 80]]
- [[_COMMUNITY_Community 81|Community 81]]
- [[_COMMUNITY_Community 82|Community 82]]
- [[_COMMUNITY_Community 83|Community 83]]
- [[_COMMUNITY_Community 84|Community 84]]
- [[_COMMUNITY_Community 85|Community 85]]
- [[_COMMUNITY_Community 86|Community 86]]
- [[_COMMUNITY_Community 87|Community 87]]
- [[_COMMUNITY_Community 88|Community 88]]
- [[_COMMUNITY_Community 89|Community 89]]
- [[_COMMUNITY_Community 90|Community 90]]
- [[_COMMUNITY_Community 91|Community 91]]
- [[_COMMUNITY_Community 92|Community 92]]

## God Nodes (most connected - your core abstractions)
1. `What You Must Do When Invoked` - 16 edges
2. `/graphify` - 15 edges
3. `API Contracts` - 14 edges
4. `Issue Validation Report` - 14 edges
5. `Frontend Patterns` - 12 edges
6. `JSON()` - 11 edges
7. `TicketRepository` - 11 edges
8. `Repository` - 11 edges
9. `compilerOptions` - 11 edges
10. `scripts` - 11 edges

## Surprising Connections (you probably didn't know these)
- `DDD Layer Separation (Domain, Repository, Service, Handler)` --semantically_similar_to--> `Backend Layered Architecture (Domain, Repository, Service, Handler)`  [INFERRED] [semantically similar]
  emergency_update_plan.md → .agents/skills/backend-go-architecture/SKILL.md
- `Technical Debt Register` --conceptually_related_to--> `Response Envelope`  [INFERRED]
  docs/tech-debt.md → .agents/skills/fullstack-api-integration/references/api-contracts.md
- `TD-014: Empty List Response Omits data Key` --semantically_similar_to--> `Important: Response omitempty Drops Empty List Data`  [INFERRED] [semantically similar]
  docs/tech-debt.md → issue.md
- `Normal Ticket Update Workflow (Forward-Only)` --semantically_similar_to--> `State Transition Side-Effect Pattern`  [INFERRED] [semantically similar]
  emergency_update_plan.md → .agents/skills/backend-go-architecture/references/backend-patterns.md
- `Important: Response omitempty Drops Empty List Data` --conceptually_related_to--> `Response Envelope`  [INFERRED]
  issue.md → .agents/skills/backend-go-architecture/references/backend-patterns.md

## Hyperedges (group relationships)
- **Ticket Domain State Machine** — emergency_update_plan_state_machine, api_openapi_ticket_schema, api_openapi_ticket_endpoints, references_backend_patterns_state_transition, issue_picked_up_reversal, issue_warranty_field_picked_up, emergency_update_plan_normal_update, emergency_update_plan_emergency_update [EXTRACTED 1.00]
- **API Contract Alignment Layer** — references_api_contracts_spec_first, references_api_contracts_generated_types, references_api_contracts_mock_api, references_api_contracts_response_envelope, api_openapi_response_envelope_schema, api_openapi_openapi_spec, api_openapi_oapi_codegen_config, references_frontend_patterns_list_discipline, references_frontend_patterns_api_payload [EXTRACTED 1.00]
- **Auth Security Chain (Cookies, RTR, Error Handling)** — changelog_auth_security_module, references_backend_patterns_auth_sessions, api_openapi_auth_endpoints, issue_token_error_leakage, issue_access_token_auto_refresh [INFERRED 0.85]

## Communities (94 total, 33 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.06
Nodes (23): main(), AccessClaims, TestMain(), AuthHandlerTestSuite, AuthRepositoryTestSuite, AuthService, CheckPasswordHash(), GenerateAccessToken() (+15 more)

### Community 1 - "Community 1"
Cohesion: 0.10
Nodes (14): PublicTrackerTicketStatus, AuthHandler, APIResponse, Error(), JSON(), maskName(), maskPhone(), ToPublicTrackerTicketAPI() (+6 more)

### Community 2 - "Community 2"
Cohesion: 0.22
Nodes (11): PublicTrackerTicket Schema (Masked PII), Ticket CRUD + Emergency Endpoints, Ticket Full Schema (Admin), Warranty Schema, Emergency Ticket Update Workflow (Reversal Capable), Normal Ticket Update Workflow (Forward-Only), Ticket Status State Machine (received, in_repair, completed, cancelled, warehouse, picked_up), Critical: picked_up Reversal Without Side-Effect Handling (+3 more)

### Community 3 - "Community 3"
Cohesion: 0.09
Nodes (18): mockAuthService, stored, KEYS, mockDbService, initMockNetwork(), initialInventory, initialSales, initialTickets (+10 more)

### Community 4 - "Community 4"
Cohesion: 0.06
Nodes (33): dependencies, lucide-svelte, devDependencies, jsdom, openapi-typescript, svelte, svelte-check, @sveltejs/adapter-node (+25 more)

### Community 5 - "Community 5"
Cohesion: 0.07
Nodes (28): ApiResponseDashboardData, ApiResponseEmpty, ApiResponseError, ApiResponseMeData, ApiResponseProductData, ApiResponseProductList, ApiResponsePublicTrackerTicketData, ApiResponseSaleData (+20 more)

### Community 6 - "Community 6"
Cohesion: 0.12
Nodes (9): Database, NewConnection(), TestDatabaseConnectionRetry(), TestMain(), DatabaseConnectionTestSuite, findMigrationPath(), SetupTestDB(), IntegrationSuite (+1 more)

### Community 8 - "Community 8"
Cohesion: 0.14
Nodes (11): PublicTrackerTicket, createdTicket, fetchSpy, mockTicket, mockTickets, newTicket, updatedTicket, updates (+3 more)

### Community 9 - "Community 9"
Cohesion: 0.12
Nodes (15): devDependencies, @playwright/test, @types/node, typescript, name, private, scripts, check (+7 more)

### Community 10 - "Community 10"
Cohesion: 0.04
Nodes (49): code:bash ($(cat graphify-out/.graphify_python) -c "), code:block11 ([Agent tool call 1: files 1-15, subagent_type="general-purpo), code:bash (PROJECT_ROOT=$(cat graphify-out/.graphify_root)), code:block13 (You are a graphify extraction subagent. Read the files liste), code:bash ($(cat graphify-out/.graphify_python) -c "), code:bash ($(cat graphify-out/.graphify_python) -c "), code:bash ($(cat graphify-out/.graphify_python) -c "), code:bash ($(cat graphify-out/.graphify_python) -c ") (+41 more)

### Community 12 - "Community 12"
Cohesion: 0.17
Nodes (9): inventoryService, Product, createdProduct, fetchSpy, mockProduct, mockProducts, newProduct, updatedProduct (+1 more)

### Community 13 - "Community 13"
Cohesion: 0.30
Nodes (3): numeric, $lib/services/inventory, targetStock

### Community 14 - "Community 14"
Cohesion: 0.33
Nodes (3): $lib/services/auth, $app/navigation, $app/state

### Community 15 - "Community 15"
Cohesion: 0.15
Nodes (12): compilerOptions, allowJs, checkJs, esModuleInterop, forceConsistentCasingInFileNames, moduleResolution, resolveJsonModule, rewriteRelativeImportExtensions (+4 more)

### Community 19 - "Community 19"
Cohesion: 0.24
Nodes (9): AppConfig, getEnv(), getEnvAsDuration(), getEnvAsInt(), LoadConfig(), TestLoadConfig_Defaults(), TestLoadConfig_PoolSettings(), TestLoadConfig_Validation() (+1 more)

### Community 20 - "Community 20"
Cohesion: 0.21
Nodes (8): adjustQty(), changeAmount, clearCart(), filteredProducts, finalTotal, handleCheckout(), loadProducts(), removeFromCart()

### Community 22 - "Community 22"
Cohesion: 0.20
Nodes (7): Sale, SaleItem, saleService, createdSale, fetchSpy, mockSales, newSale

### Community 23 - "Community 23"
Cohesion: 0.20
Nodes (9): compilerOptions, esModuleInterop, forceConsistentCasingInFileNames, module, moduleResolution, skipLibCheck, strict, target (+1 more)

### Community 25 - "Community 25"
Cohesion: 0.22
Nodes (8): emailInput, errorAlert, heading, logoutButton, passwordInput, submitButton, url, welcomeHeader

### Community 27 - "Community 27"
Cohesion: 0.08
Nodes (12): TicketUpdateDevicePosition, TicketUpdatePaymentMethod, TicketUpdateStatus, AdminTicketService, NewAdminService(), setupServiceTest(), TestService_CreateTicket(), TestService_GetTicket() (+4 more)

### Community 28 - "Community 28"
Cohesion: 0.25
Nodes (5): apiFetch(), fetchSpy, mockWarranties, Warranty, warrantyService

### Community 29 - "Community 29"
Cohesion: 0.25
Nodes (5): $lib/services/warranty, daysLeft, expiredWarrantiesCount, filteredWarranties, status

### Community 30 - "Community 30"
Cohesion: 0.05
Nodes (38): code:block1 (/graphify                                             # full), code:bash (if [ ! -f graphify-out/.graphify_python ]; then), code:bash ($(cat graphify-out/.graphify_python) -c "), code:bash ($(cat graphify-out/.graphify_python) -c "), code:bash ($(cat graphify-out/.graphify_python) -c "), code:bash (if [ ! -f graphify-out/.graphify_extract.json ]; then), code:bash ($(cat graphify-out/.graphify_python) -c "), code:bash ($(cat graphify-out/.graphify_python) -c ") (+30 more)

### Community 31 - "Community 31"
Cohesion: 0.29
Nodes (6): card, catalogItem, detailUrl, searchInput, tId, url

### Community 32 - "Community 32"
Cohesion: 0.33
Nodes (6): button, button2, createTicketAndGetId(), input, input2, url

### Community 34 - "Community 34"
Cohesion: 0.18
Nodes (6): button, child, childrenSnippet, { getByRole }, { getByTestId }, handleClick

### Community 35 - "Community 35"
Cohesion: 0.09
Nodes (21): 10. Ambiguous Status Mapping (`ready_for_pickup`), 11. Startup Refresh Token Purge (`TD-003`), 12. Duplicate Mapping Logic, 13. Public Tracker Routing Conflict, 1. Unit Test Handler Cookie Secure Behavior, 2. Refresh Token Route Mapping in Tests, 3. Ticket ID Format Validation, 4. Empty List Response Omitting `data` Key (`TD-014`) (+13 more)

### Community 36 - "Community 36"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 37 - "Community 37"
Cohesion: 0.33
Nodes (5): MeResponse, RefreshResult, SignInRequest, SignInResponse, SignInResult

### Community 40 - "Community 40"
Cohesion: 0.12
Nodes (16): 1. Daftar Kemungkinan Data yang Dapat Diubah, 2. Rules (Invariant) yang Tetap Berlaku pada Emergency Update, 3. Rencana Pemisahan (Backend Plan), 4. Rencana Pemisahan UI (Frontend Svelte 5 Plan), 5. Rencana UI Quick Actions (Halaman Daftar Tiket `/admin/tickets`), 6. Rencana Kerja Urutan Prioritas (Modul A: Normalisasi Database Terlebih Dahulu), A. Domain Layer (`apps/backend/internal/ticket/model.go`), A. Normal Update (Update Tiket Biasa) (+8 more)

### Community 42 - "Community 42"
Cohesion: 0.25
Nodes (5): ../components/ProductForm.svelte, numeric, displayCostPrice, displayPrice, $lib/utils/format

### Community 43 - "Community 43"
Cohesion: 0.50
Nodes (4): Conventional Commits Standard, Emergency Update & Normalization Plan, Code Review Verdict: Not Ready for Merge, Code Review & Ops Checklist Framework

### Community 44 - "Community 44"
Cohesion: 0.67
Nodes (3): OpenBench Go/Fiber + SvelteKit + Playwright Monorepo, PostgreSQL Dev Database Service, Test Infrastructure Stack

### Community 51 - "Community 51"
Cohesion: 0.17
Nodes (11): E2E Test Package Separation from Frontend, Backend Commands, Backend Patterns, code:typescript (await page.waitForSelector('main[data-hydrated="true"]');), Container And Compose Commands, Frontend Commands, Frontend Patterns, Playwright Patterns (+3 more)

### Community 59 - "Community 59"
Cohesion: 0.12
Nodes (15): API Payload Discipline, Async Views, code:typescript (class ToastService {), code:svelte (let rawAmount = $state(0);), code:typescript (let fetchId = 0;), Component Slicing, Derived API Data In UI, Files To Inspect First (+7 more)

### Community 67 - "Community 67"
Cohesion: 0.15
Nodes (12): API Contracts, Auth And Cookies, Contract Source, Files To Inspect First, Generated Type Gate, Generator Toolchain Compatibility, Mock Toggle Reconciliation, Naming Rules (+4 more)

### Community 68 - "Community 68"
Cohesion: 0.20
Nodes (9): Agent Skills & Guidelines, Build, Test, and Development Commands, Coding Style & Naming Conventions, Commit & Pull Request Guidelines, Knowledge Graph (graphify), Project Structure & Module Organization, Repository Guidelines, Security & Configuration Tips (+1 more)

### Community 69 - "Community 69"
Cohesion: 0.20
Nodes (9): Backend Patterns, Config And Server Safety, Files To Inspect First, Layer Responsibilities, Migrations, Public Endpoints, Sequential Identifiers, State Transitions (+1 more)

### Community 70 - "Community 70"
Cohesion: 0.20
Nodes (9): Backend, Directory Conventions, Frontend & UI, Fullstack & Integration, Merge History, OpenBench Skills Catalog, Skills Directory, Testing (+1 more)

### Community 71 - "Community 71"
Cohesion: 0.20
Nodes (9): Active Navigation, Cards And Surfaces, Dashboards, Files To Inspect First, Loading States, Responsive Checks, Theme Tokens, UI Patterns (+1 more)

### Community 72 - "Community 72"
Cohesion: 0.22
Nodes (8): Added, Changed, Changelog, Changelog Standard, Fixed, Removed, Security, [Unreleased]

### Community 73 - "Community 73"
Cohesion: 0.25
Nodes (7): Code Review, Commit Readiness, Dependency And Toolchain, Docker And Compose, Environment And Config, Review And Ops Checklist, Technical Debt

### Community 74 - "Community 74"
Cohesion: 0.38
Nodes (7): oapi-codegen Configuration, OpenBench OpenAPI 3.0.3 Specification, Health Probe Module (Liveness, Readiness), DDD Layer Separation (Domain, Repository, Service, Handler), Generated Type Gate, Spec-First Workflow (OpenAPI as Single Source of Truth), API Payload Discipline (Generated Request Types)

### Community 75 - "Community 75"
Cohesion: 0.38
Nodes (7): API Response Envelope {code, message, data}, TD-014: Empty List Response Omits data Key, Important: Response omitempty Drops Empty List Data, code:json ({), Mock API, Response Envelope, List Data Discipline (Empty Array Contract)

### Community 76 - "Community 76"
Cohesion: 0.33
Nodes (6): Curated Agent Skills System, Backend Layered Architecture (Domain, Repository, Service, Handler), code:go (return response.JSON(c, fiber.StatusOK, "Message", payload)), Response Envelope, Service-Layer Transaction Pattern, Six Core Agent Skills Catalog

### Community 77 - "Community 77"
Cohesion: 0.33
Nodes (5): Backend Go Architecture, Hard Checks, Load References, Operating Rule, Workflow

### Community 78 - "Community 78"
Cohesion: 0.33
Nodes (5): Frontend Svelte 5 Architecture, Hard Checks, Load References, Operating Rule, Workflow

### Community 79 - "Community 79"
Cohesion: 0.33
Nodes (5): Fullstack API Integration & Mocking, Hard Checks, Load References, Operating Rule, Workflow

### Community 80 - "Community 80"
Cohesion: 0.33
Nodes (5): Hard Checks, Load References, OpenBench Testing Strategy, Operating Rule, Workflow

### Community 81 - "Community 81"
Cohesion: 0.33
Nodes (5): Hard Checks, Load References, OpenBench UI Design System, Operating Rule, Workflow

### Community 82 - "Community 82"
Cohesion: 0.33
Nodes (5): Hard Checks, Load References, OpenBench Workflow & Ops, Operating Rule, Workflow

### Community 83 - "Community 83"
Cohesion: 0.40
Nodes (6): Auth Cookie-Based Endpoints (signin, refresh, signout, me), Auth Security Module (RTR, Bcrypt, HttpOnly Cookies), Technical Debt Register, Missing Auto-Refresh on Access Token Expiry in Ticket Detail View, Important: Refresh Token Error Detail Leakage, Auth Sessions

### Community 84 - "Community 84"
Cohesion: 0.33
Nodes (5): $lib/components/ToastContainer.svelte, $lib/assets/favicon.svg, $lib/services/toast.svelte, ./layout.css, svelte/transition

### Community 85 - "Community 85"
Cohesion: 0.20
Nodes (8): authService, stored, fetchMock, mockResponse, mockSignInResponse, mockSignOutResponse, stored, UserSession

## Knowledge Gaps
- **380 isolated node(s):** `AppConfig`, `AccessClaims`, `SignInRequest`, `SignInResult`, `RefreshResult` (+375 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **33 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `main()` connect `Community 0` to `Community 27`, `Community 19`, `Community 6`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Why does `NewAdminService()` connect `Community 27` to `Community 0`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **What connects `AppConfig`, `AccessClaims`, `SignInRequest` to the rest of the system?**
  _382 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 0` be split into smaller, more focused modules?**
  _Cohesion score 0.06292517006802721 - nodes in this community are weakly interconnected._
- **Should `Community 1` be split into smaller, more focused modules?**
  _Cohesion score 0.09803921568627451 - nodes in this community are weakly interconnected._
- **Should `Community 3` be split into smaller, more focused modules?**
  _Cohesion score 0.09009009009009009 - nodes in this community are weakly interconnected._
- **Should `Community 4` be split into smaller, more focused modules?**
  _Cohesion score 0.058823529411764705 - nodes in this community are weakly interconnected._