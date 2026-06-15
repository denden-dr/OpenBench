---
name: openbench-workflow-and-ops
description: Use when setting up dev environments, Docker/Podman compose, securing container builds, committing code (Atomic Commits), or performing code reviews.
version: 1.0.0
---

# OpenBench Workflow & Ops

## Overview
Guidelines for local development environments, container security, atomic Git commits, and code review audits.

## Dev Environment & Compose
1. **Docker Compose**: Manage PostgreSQL via Compose with proper health checks (`pg_isready`). Include Makefile teardown targets.
2. **Config Loader**: Load environment variables (`.env`, `.env.test`) securely based on `APP_ENV`. Implement production safety gates to reject empty database passwords or `DB_SSLMODE=disable` outside local environments.
3. **Go Connection Pooling**: Configure `sqlx` pooling with explicit context timeouts during retries. Expose `db.DB.Stats()` for observability.

## Secure Containerization
1. **`.dockerignore`**: Exclude `.env`, `.git`, and `node_modules` from the context.
2. **Multi-Stage Builds**: Use multi-stage Dockerfiles for Go backend (static binaries) and Node/SvelteKit frontend (building assets vs running).
3. **Non-Root Execution**: Use `USER appuser` or `USER node` at the end of the Dockerfile to prevent root-breakout vulnerabilities. Use pinned image tags (e.g., `golang:1.23-alpine` instead of `latest`).

## Practicing Atomic Commits
1. **Single Logical Change**: A commit must represent a single, compilable change with passing tests.
2. **Interactive Staging**: Use `git add -p` to stage specific hunks and avoid "while-I'm-here" commits.
3. **Conventional Commits**: Format messages as `type(scope): description`.
4. **Check Before Commit**: Do not blindly `git add .`. Do not mix refactoring and feature additions in a single commit.

## Code Review Workflows
1. **Feature-Scoped or Full Audit**: Always review code across Architecture, Security, Testing, Best-Practices, and Performance.
2. **Reference Checklist**: Adhere to `references/review-checklist.md`.
3. **Track Technical Debt**: Log deferrable compromises in `docs/tech-debt.md`. Document active bugs in `hotissue.md`.
4. **Key Verification Areas**:
   - No hardcoded secrets.
   - Svelte 5 runes properly utilized; dead code eliminated.
   - DB connection pool limits respected.

## Common Mistakes to Avoid
- Storing global DB connections without stats observability.
- Committing sensitive `.env` files.
- Running containers as root.
- Baking local `node_modules` into container images.
- Committing broken code or using commit messages with "and" (`feat A and B`).
