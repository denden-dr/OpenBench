---
name: openbench-workflow-and-ops
description: Operate the OpenBench monorepo safely. Use when setting up local dev, Docker/Podman compose, environment config, container builds, package-manager installs, Makefile workflows, commits, pull-request readiness, code reviews, security checks, or technical-debt tracking.
---

# OpenBench Workflow & Ops

## Operating Rule

Prefer existing Makefile and compose targets. Keep environment, container, and commit changes explicit and scoped.

## Workflow

1. Inspect the Makefile and compose files before inventing commands.
2. For environment changes, update examples and safety gates without committing real secrets.
3. For container or dependency-install changes, inspect Dockerfiles, package manifests, lockfiles, and package-manager strict install behavior together.
4. For reviews, lead with concrete findings and file/line references.
5. For commits, stage only the logical change and use Conventional Commits.
6. Record active blocking bugs in root `issue.md`; record deliberate, deferrable compromises in `docs/tech-debt.md` when that file exists.

## Load References

- Read `references/review-checklist.md` for review, commit readiness, environment/config, Docker, Makefile, or release-prep work.

## Hard Checks

- Do not commit `.env` files or secrets.
- Do not use broad `git add .` when unrelated work exists.
- Do not run containers as root in final production images.
- Do not use unpinned `latest` base images.
- Do not use package-manager bypass flags such as `--force` or `--legacy-peer-deps` as a permanent fix for peer dependency conflicts.
- Keep package manifests and lockfiles in sync when dependency versions change.
