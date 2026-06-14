---
name: requesting-code-review
description: Use when performing code reviews in the repository, either scoped to feature diffs or across the entire codebase to verify monorepo health, security, testing robustness, and config alignment.
version: 1.2.0
---

# Requesting Code Review

## Overview
Code reviews verify that changes are clean, safe, robust, and aligned with monorepo patterns. Reviews are scoped as Feature Code Reviews (scoped to branch diff) or Full Codebase Reviews (pre-release health check).

## When to Use
- Auditing local branch diffs before completing feature development.
- Checking codebase alignment across packages before production releases.
- Evaluating changes against Monorepo architecture constraints.

### When NOT to Use
- General programming guidance unrelated to code audits or project review workflows.

---

## Core Process: Multi-Dimensional Review

Every code review must evaluate modifications across five dimensions:
1.  **Architecture & Boundaries** (Monorepo boundaries, Go domain isolation, dependency cycles)
2.  **Security** (No hardcoded secrets, cookie flags, token rotation, input validation)
3.  **Testing** (Port isolation, container cleanup, mock assertions, hydration delays)
4.  **Best-Practices** (Structured logging, Svelte 5 runes, dead code elimination, conventional commits)
5.  **Performance** (DB connection pooling, Vite module pre-bundling, Playwright concurrency limit)

> [!IMPORTANT]
> **REQUIRED REFERENCE:** Always read and follow the specific audit checks in [review-checklist.md](./references/review-checklist.md) when evaluating code changes.

---

## Step-by-Step Instructions

1.  **Determine Review Scope**: Select either a Feature-Scoped Review (auditing only the current branch diff) or a Full Codebase Review.
2.  **Execute Feature-Scoped Audit**:
    *   Run `git diff --name-only main...HEAD` to isolate modified files.
    *   Audit each file against the checklist in [review-checklist.md](./references/review-checklist.md).
    *   Run unit and integration tests targeting the modified packages.
3.  **Execute Full Codebase Audit**:
    *   Verify environment configurations in `.env` match variables loaded in code.
    *   Run the full monorepo test suites (`go test ./...` in backend, `npm run check` in frontend).
4.  **Log and Track Findings**:
    *   Register active, blocking bugs and security risks in `hotissue.md` at the repository root.
    *   Register standard improvements and fixes in `issue.md` at the repository root.
    *   Document deferrable compromises and technical debt in `docs/tech-debt.md` using the format in [tech-debt-template.md](./references/tech-debt-template.md). Classify entries as Architecture, Database, Testing, Security, or Performance.

---

## Quick Reference: Review Scopes

| Review Dimension | Feature-Scoped Audit | Full Codebase Audit |
|---|---|---|
| **Architecture** | Import boundary checks | Directory layout & package boundary checks |
| **Security** | Secret checks & input validation | CORS configurations & global middleware |
| **Testing** | Dynamic test setup & local pass checks | Monorepo CI verification & mock coverage |
| **Best-Practice** | Style, error wrapping, Svelte runes | Conventional commits, contract consistency |
| **Performance** | Goroutine leaks, Vite optimizer | DB pool limits, worker parallelism tuning |
