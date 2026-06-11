# Technical Debt Register

This document tracks compromises, shortcuts, or design issues in the OpenBench codebase that require remediation.

| ID | Category | Title | Description / Context | Impact | Remediation Plan | Effort | Status |
|----|----------|-------|-----------------------|--------|------------------|--------|--------|
| TD-001 | Architecture | Implicit Env File Discovery | `apps/backend/internal/config/config.go` still selects `.env` files from the current directory, `..`, or `../..`. This is narrower than the original parent-tree walk, but startup behavior still depends on the process working directory. | Local, CI, and container runs can load different env files unless launched from expected directories. | Resolve env files from an explicit path or fixed project-root convention. | Medium | Active |
| TD-002 | Performance | Missing Database Pool Observability | `apps/backend/internal/database/database.go` exposes `Stats()`, but no logs, metrics endpoint, or alerting path consumes pool stats or slow-operation signals yet. | Latency regressions will still be hard to distinguish between pool starvation and slow queries. | Emit `db.Stats()` through logs/metrics and add slow query or connection wait instrumentation. | Medium | Active |
