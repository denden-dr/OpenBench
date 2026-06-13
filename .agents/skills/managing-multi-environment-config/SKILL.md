---
name: managing-multi-environment-config
description: Use when setting up environment-based configurations in a Go application, managing different environments (development, testing, production) using .env files, and loading variables cleanly. Do not use for frontend environment configuration or non-Go projects.
version: 1.0.0
---

# Managing Multi-Environment Config

## Overview
Go applications should support running in different environments (development, testing, production) seamlessly. The configuration should be loaded dynamically based on the current environment, controlled by a standardized `APP_ENV` environment variable, ensuring that test runs do not affect development databases or systems.

## When to Use
- Designing the configuration loader for a Go service.
- Defining environment files (`.env`, `.env.test`, `.env.example`).
- Isolating development and testing parameters (e.g., database host, port, credentials).

## Step-by-Step Instructions

1. **Verify Environment Variables**: Use `assets/env.example.template` to structure local `.env`, `.env.test`, and `.env.example` templates, keeping sensitive secrets out of source control.
2. **Implement Config Struct & Loader**: Read `assets/config.go.template` and construct the configurations loading script, ensuring it reads `APP_ENV` to determine whether to load `.env` or `.env.test`. Normalize values like `"dev"` to `"development"` and validate that `APP_ENV` strictly matches one of the expected values (`"development"`, `"test"`, `"production"`).
3. **Handle Directory Discoverability**: Restrict `.env` file discovery to the local working directory or at most one parent folder level up to prevent infinite upward traversal.
4. **Implement Production Safety Gates**: Add checks in the configuration loader to reject starting the application if the environment is non-local (not "development" or "test") and database credentials are weak or SSL mode is disabled.
5. **Centralize Fallback Constants**: Ensure default fallback constants (like CORS origins or timeouts) are defined exclusively inside the configuration loader struct, eliminating redundant fallback copies in router initialization paths.

## Common Mistakes
- **Implicit Env File Discovery (TD-003)**: Walking up parent directories indefinitely to find `.env` files. This makes configuration load-dependent on the current working directory. Keep it to a single level or explicit path parameter.
- **Hardcoded CORS Origins (TD-002)**: Hardcoding allowed origins inside the router or application startup code. Always pull allowed origins from the configuration.
- **ENV/APP_ENV Inconsistency**: Mixing up variable names like `ENV` in docker-compose or .env files but reading `APP_ENV` in the application loader (BE-003).
- **No Production Safety Gates**: Allowing the application to start up in a production environment with insecure local defaults like an empty database password or `DB_SSLMODE=disable` (BE-004).
- **Committing Sensitive Secrets**: Committing actual `.env` files with production or local personal passwords to Git. Always add `.env` to `.gitignore`.
