---
name: documenting-technical-debt
description: Use when documenting workarounds, legacy shortcuts, or architectural compromises in the codebase, ensuring they are cataloged systematically in docs/tech-debt.md. Do not use for standard code refactoring that does not represent technical debt.
version: 1.0.0
---

# Documenting Technical Debt

## Overview
Technical debt is an inevitable part of software development. However, undocumented technical debt is dangerous because it gets forgotten and leads to maintenance overhead, bugs, or security vulnerabilities. We track all technical debt systematically in `/docs/tech-debt.md`.

## When to Use
- Whenever you introduce a temporary workaround or shortcut due to time limits, dependency issues, or environment constraints.
- When you identify a suboptimal design pattern in existing code that should be refactored.
- Before completing feature development, if you had to compromise on best practices.

## Step-by-Step Instructions

1. **Identify and Classify**: Identify the technical compromise and categorize it under one of the standard classifications: Architecture, Database, Testing, Security, or Performance.
2. **Review the Entry Template**: Read `assets/tech-debt-template.md` to see the structure and formatting of the technical debt register table.
3. **Record the Entry**: Append a new row to the table in `/docs/tech-debt.md`. Populate the ID, Category, Title, Description/Context, Impact/Risk, Remediation Plan, Effort (Low/Medium/High), and Status fields. Make sure to link back to the affected source files/lines.

## Common Mistakes
- **Vague Descriptions**: Writing "refactor backend code" without stating what is wrong and where.
- **No Remediation Plan**: Recording a problem without detailing the technical steps to solve it.
- **Not Linking Source**: Forgetting to link to the specific file or package containing the debt.
