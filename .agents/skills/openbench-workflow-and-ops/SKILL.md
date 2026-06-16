---
name: openbench-workflow-and-ops
description: Operate the OpenBench monorepo safely. Use when setting up local dev, Docker/Podman compose, environment config, container builds, Makefile workflows, commits, pull-request readiness, code reviews, security checks, or technical-debt tracking.
---

# OpenBench Workflow & Ops

## Operating Rule

Prefer existing Makefile and compose targets. Keep environment, container, and commit changes explicit and scoped.

## Workflow

1. Inspect the Makefile and compose files before inventing commands.
2. For environment changes, update examples and safety gates without committing real secrets.
3. For container changes, keep builds multi-stage, pinned, and non-root.
4. For reviews, lead with concrete findings and file/line references.
5. For commits, stage only the logical change and use Conventional Commits.
6. Record deliberate, deferrable compromises in `docs/tech-debt.md` when that file exists; record active blocking bugs in `hotissue.md` when that convention exists.

## Load References

- Read `references/review-checklist.md` for review, commit readiness, environment/config, Docker, Makefile, or release-prep work.

## Hard Checks

- Do not commit `.env` files or secrets.
- Do not use broad `git add .` when unrelated work exists.
- Do not run containers as root in final production images.
- Do not use unpinned `latest` base images.
