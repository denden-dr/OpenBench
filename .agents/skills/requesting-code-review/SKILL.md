---
name: requesting-code-review
description: Use when performing code reviews in the repository, either scoped to feature diffs or across the entire codebase to verify monorepo health and config alignment. Do not use for external code reviews or non-monorepo architectures.
version: 1.0.0
---

# Requesting Code Review

## Overview
Code reviews verify cleanliness, safety, database connection health, security, and Monorepo config alignment. Depending on the size of the changes, reviews are scoped in two ways:
1. **Feature Code Review**: Focused only on the current feature branch's diff.
2. **Full Codebase Review**: Checking directory architecture, global variables, and cross-application interactions before a release.

## When to Use
- Auditing local branch diffs before completing feature development.
- Checking codebase alignment across packages before production releases.
- Reorganizing directory layouts or updating global env variables.

## Step-by-Step Instructions

1. **Determine Review Scope**: Select either a Feature-Scoped Review (auditing only the current branch diff) or a Full Codebase Review (auditing the entire project structure).
2. **Execute Feature-Scoped Audit**:
   - Run `git diff --name-only main...HEAD` to isolate modified files.
   - Audit each file for code styling, database resource/pool handling, security middlewares, and credential leaks.
   - Run unit tests targeting the modified packages.
3. **Execute Full Codebase Audit**:
   - Inspect the monorepo directory layout to ensure apps are isolated under `/apps`.
   - Verify environment configurations in `.env` match variables loaded in code and docker compose setups.
   - Run the full monorepo test suites (`go test ./...` in backend, `npm run check` in frontend).
4. **Log and Track Findings**:
   - Register active, blocking bugs and security risks in `hotissue.md` at the repository root.
   - Register standard improvements and fixes in `issue.md` at the repository root.
   - Document deferrable compromises and technical debt in `docs/tech-debt.md`.
