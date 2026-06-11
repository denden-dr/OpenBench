---
name: requesting-feature-code-review
description: Use when initiating a code review for changes introduced in the current feature branch, comparing changes against the base branch.
---

# Requesting Feature Code Review

## Overview
Feature-scoped code review focuses only on the changes introduced by the current branch compared to the base branch (usually `main`). This limits the cognitive load of the review and ensures that only modified files, new features, or bug fixes are checked for regressions, formatting, and performance.

## When to Use
- When feature branch development is complete and you want to review only the diff before merging.
- To check if the changes adhere to project coding guidelines without reviewing unaffected parts of the codebase.

## Core Process

### 1. Identify Modified Files
Run git commands to find files modified in the feature branch relative to the base branch (e.g. `main`):
```bash
git diff --name-only main...HEAD
```

### 2. Perform Diff Audit
For each modified file, review:
- **Cleanliness & Formatting**: Check if new imports are used, and formatting follows idiomatic Go/JS conventions.
- **Connection Pools & Resources**: Ensure any new database queries or connections use the injected `*sqlx.DB` instance and do not leak resources.
- **Middlewares & Security**: Ensure new endpoints are secure and covered by existing middlewares.
- **No Hardcoded Credentials**: Confirm no credentials or absolute hostnames are added.

### 3. Verify Scoped Tests
Run tests specifically targeting the packages or files that were modified:
```bash
go test -v ./internal/database
```

## Response Patterns
Present findings specifically highlighting:
- Files modified.
- Technical feedback on the diff.
- Execution results of the related tests.
