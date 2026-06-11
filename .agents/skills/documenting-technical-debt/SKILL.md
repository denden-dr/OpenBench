---
name: documenting-technical-debt
description: Use when documenting workarounds, legacy shortcuts, or architectural compromises in the codebase, ensuring they are cataloged systematically in docs/tech-debt.md.
---

# Documenting Technical Debt

## Overview
Technical debt is an inevitable part of software development. However, undocumented technical debt is dangerous because it gets forgotten and leads to maintenance overhead, bugs, or security vulnerabilities. We track all technical debt systematically in `/docs/tech-debt.md`.

## When to Use
- Whenever you introduce a temporary workaround or shortcut due to time limits, dependency issues, or environment constraints.
- When you identify a suboptimal design pattern in existing code that should be refactored.
- Before completing feature development, if you had to compromise on best practices.

## Core Process

### 1. Identify and Categorize
Classify the debt into one of these categories:
- **Architecture**: Suboptimal module boundaries, tight coupling.
- **Database**: Poor schema design, inefficient queries, index gaps.
- **Testing**: Low test coverage, flaky tests, lack of integration tests.
- **Security**: Minor security compromises (e.g., permissive CORS in local dev, missing inputs validation to be added later).
- **Performance**: High memory utilization, missing cache.

### 2. Record the Entry
Add a new entry to the bottom of the active table in `/docs/tech-debt.md` with:
- **Title**: Descriptive title.
- **Description**: What is the compromise, and *why* was it introduced (context).
- **Impact/Risk**: What happens if it remains unresolved (e.g., performance degradation, security leak, developer velocity slowing down).
- **Remediation Plan**: Actionable steps to resolve the debt.
- **Effort**: Estimate (Low / Medium / High).
- **File Reference**: Links to the source files/lines.

## The Standard Entry Template (docs/tech-debt.md)

Entries in `docs/tech-debt.md` should use the following markdown table format:

| ID | Category | Title | Description / Context | Impact | Remediation Plan | Effort |  Status |
|----|----------|-------|-----------------------|--------|------------------|--------|--------|
| TD-001 | Database | Ephemeral Test Volumes | Postgres-test runs on ephemeral storage without volume mounts. | Test data is wiped on container down. Cannot debug persistent issues between runs. | Set up separate persistent test volume if state persistence becomes necessary. | Low | Active |

## Common Mistakes
- **Vague Descriptions**: Writing "refactor backend code" without stating what is wrong and where.
- **No Remediation Plan**: Recording a problem without detailing the technical steps to solve it.
- **Not Linking Source**: Forgetting to link to the specific file or package containing the debt.
