# Graph Report - OpenBench  (2026-05-24)

## Corpus Check
- 27 files · ~79,896 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 181 nodes · 229 edges · 21 communities detected
- Extraction: 74% EXTRACTED · 26% INFERRED · 0% AMBIGUOUS · INFERRED: 59 edges (avg confidence: 0.83)
- Token cost: 0 input · 0 output

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
- [[_COMMUNITY_Community 17|Community 17]]
- [[_COMMUNITY_Community 25|Community 25]]
- [[_COMMUNITY_Community 26|Community 26]]
- [[_COMMUNITY_Community 27|Community 27]]
- [[_COMMUNITY_Community 28|Community 28]]

## God Nodes (most connected - your core abstractions)
1. `Ticket` - 12 edges
2. `TicketIntegrationTestSuite` - 10 edges
3. `MapDatabaseError()` - 8 edges
4. `NewTicketService()` - 8 edges
5. `MapRepositoryError()` - 8 edges
6. `main()` - 7 edges
7. `runTests()` - 7 edges
8. `TicketService` - 7 edges
9. `MockTicketRepository` - 7 edges
10. `TicketHandler` - 6 edges

## Surprising Connections (you probably didn't know these)
- `PhoneFix Admin Overview` --semantically_similar_to--> `Phone Repair PRD`  [INFERRED] [semantically similar]
  README.md → docs/PRD.md
- `PhoneFix Admin Overview` --conceptually_related_to--> `Local Compose Stack`  [INFERRED]
  README.md → compose.yaml
- `main()` --calls--> `Load()`  [INFERRED]
  backend/main.go → backend/internal/config/config.go
- `main()` --calls--> `NewTicketRepository()`  [INFERRED]
  backend/main.go → backend/internal/repository/ticket_repo.go
- `main()` --calls--> `NewTicketHandler()`  [INFERRED]
  backend/main.go → backend/internal/handler/ticket_handler.go

## Hyperedges (group relationships)
- **Svelte Project Scaffold Artifacts** — frontend_scaffold_readme, sveltekit_adapter_node_config, svelte_app_namespace_declarations, lib_alias_placeholder [INFERRED 0.82]
- **Frontend API Proxy Surface** — frontend_api_proxy_vite, frontend_api_proxy_hook, backend_url_environment_variable, api_proxy_pipeline [INFERRED 0.92]
- **Repair Management UI Flow** — repair_ticket_dashboard, repair_ticket_statistics, repair_ticket_create_modal, repair_ticket_edit_drawer, ticket_crud_flow [INFERRED 0.94]
- **Phone Repair Bench Setup** — equipment_inspection_microscope, equipment_hot_air_rework_station, equipment_precision_screwdriver_set, equipment_fine_tweezers [INFERRED 0.93]
- **Phone Board Microsoldering Workflow** — microsoldering_repair_scene, smartphone_logic_board, soldering_iron, precision_tweezers, a15_bionic_chip [INFERRED 0.92]
- **Ticket Lifecycle Workflow** — phone_repair_prd, ticket_model, ticket_service, ticket_handler_integration_suite [INFERRED 0.88]

## Communities

### Community 0 - "Community 0"
Cohesion: 0.14
Nodes (9): TestTicketErrorsAreStable(), TestValidateTicketUpdate(), ValidateTicketUpdate(), AppError, MapModelError(), MapRepositoryError(), TestMapModelError(), TestMapRepositoryError() (+1 more)

### Community 1 - "Community 1"
Cohesion: 0.15
Nodes (5): MockTicketRepository_Create_Call, MockTicketRepository_Delete_Call, MockTicketRepository_GetByID_Call, MockTicketRepository_List_Call, MockTicketRepository_Update_Call

### Community 2 - "Community 2"
Cohesion: 0.2
Nodes (5): Ticket, optionalText(), TicketPaymentStatus, TicketStatus, TicketUpdate

### Community 3 - "Community 3"
Cohesion: 0.22
Nodes (9): main(), NewDB(), ErrorHandler(), NewTicketService(), TestTicketService_CreateTicket(), TestTicketService_DeleteTicket(), TestTicketService_GetTicket(), TestTicketService_ListTickets() (+1 more)

### Community 4 - "Community 4"
Cohesion: 0.24
Nodes (12): Config, DefaultDatabaseConfig(), getEnv(), getEnvDuration(), getEnvInt(), Load(), loadDatabaseConfig(), DatabaseConfig (+4 more)

### Community 5 - "Community 5"
Cohesion: 0.15
Nodes (5): SetupTestDB(), NewTicketHandler(), TicketHandler, NewTicketRepository(), TicketRepository

### Community 6 - "Community 6"
Cohesion: 0.17
Nodes (3): CleanTestDB(), TestTicketIntegrationSuite(), TicketIntegrationTestSuite

### Community 7 - "Community 7"
Cohesion: 0.29
Nodes (3): MapDatabaseError(), TestMapDatabaseError(), sqlTicketRepository

### Community 8 - "Community 8"
Cohesion: 0.29
Nodes (1): MockTicketRepository

### Community 9 - "Community 9"
Cohesion: 0.29
Nodes (7): Electronics Repair Workstation, Fine Tweezers, Hot Air Rework Station, Inspection Microscope, Logic Board, Precision Screwdriver Set, Repair Mat

### Community 10 - "Community 10"
Cohesion: 0.33
Nodes (1): MockTicketRepository_Expecter

### Community 11 - "Community 11"
Cohesion: 0.33
Nodes (6): A15 Bionic Chip, Hero Image, Microsoldering Repair Scene, Precision Tweezers, Smartphone Logic Board, Soldering Iron

### Community 12 - "Community 12"
Cohesion: 0.5
Nodes (3): CreateTicketRequest, TicketResponse, UpdateTicketRequest

### Community 13 - "Community 13"
Cohesion: 0.5
Nodes (1): Error

### Community 14 - "Community 14"
Cohesion: 0.67
Nodes (3): OpenBench Favicon, Repair Tools, Smartphone Outline

### Community 15 - "Community 15"
Cohesion: 0.67
Nodes (3): Local Compose Stack, Phone Repair PRD, PhoneFix Admin Overview

### Community 17 - "Community 17"
Cohesion: 1.0
Nodes (2): Tailwind Build Failure Log, Unknown Tailwind Utility bg-background

### Community 25 - "Community 25"
Cohesion: 1.0
Nodes (1): Svelte Scaffold README

### Community 26 - "Community 26"
Cohesion: 1.0
Nodes (1): Crawler Policy

### Community 27 - "Community 27"
Cohesion: 1.0
Nodes (1): Svelte Logo

### Community 28 - "Community 28"
Cohesion: 1.0
Nodes (1): Test Compose Stack

## Knowledge Gaps
- **28 isolated node(s):** `Config`, `DatabaseConfig`, `CreateTicketRequest`, `UpdateTicketRequest`, `TicketResponse` (+23 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Community 8`** (7 nodes): `MockTicketRepository`, `.Create()`, `.Delete()`, `.EXPECT()`, `.GetByID()`, `.List()`, `.Update()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 10`** (6 nodes): `MockTicketRepository_Expecter`, `.Create()`, `.Delete()`, `.GetByID()`, `.List()`, `.Update()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 13`** (4 nodes): `errors.go`, `Error`, `.Error()`, `NewModelError()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 17`** (2 nodes): `Tailwind Build Failure Log`, `Unknown Tailwind Utility bg-background`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 25`** (1 nodes): `Svelte Scaffold README`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 26`** (1 nodes): `Crawler Policy`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 27`** (1 nodes): `Svelte Logo`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 28`** (1 nodes): `Test Compose Stack`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `main()` connect `Community 3` to `Community 4`, `Community 5`?**
  _High betweenness centrality (0.267) - this node is a cross-community bridge._
- **Why does `runTests()` connect `Community 4` to `Community 1`, `Community 3`?**
  _High betweenness centrality (0.219) - this node is a cross-community bridge._
- **Why does `NewDB()` connect `Community 3` to `Community 4`?**
  _High betweenness centrality (0.187) - this node is a cross-community bridge._
- **Are the 7 inferred relationships involving `MapDatabaseError()` (e.g. with `.Is()` and `.Create()`) actually correct?**
  _`MapDatabaseError()` has 7 INFERRED edges - model-reasoned connections that need verification._
- **Are the 7 inferred relationships involving `NewTicketService()` (e.g. with `main()` and `.SetupSuite()`) actually correct?**
  _`NewTicketService()` has 7 INFERRED edges - model-reasoned connections that need verification._
- **Are the 6 inferred relationships involving `MapRepositoryError()` (e.g. with `.CreateTicket()` and `.GetTicket()`) actually correct?**
  _`MapRepositoryError()` has 6 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Config`, `DatabaseConfig`, `CreateTicketRequest` to the rest of the system?**
  _28 weakly-connected nodes found - possible documentation gaps or missing edges._