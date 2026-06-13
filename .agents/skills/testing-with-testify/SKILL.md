---
name: testing-with-testify
description: Use when writing or structuring Go unit tests, creating mocks with testify, using assert vs require, or setting up test suites with stretchr/testify. Do not use for frontend Vitest/Jest unit tests or non-Go unit testing libraries.
version: 1.0.0
---

# Testing with Testify

## Overview
Go's built-in testing library is powerful, but writing assertions manually can lead to verbose boilerplate. The `github.com/stretchr/testify` library provides friendly assertion functions, mock helpers, and test suite structures that make unit and integration tests cleaner, more readable, and easier to maintain.

## When to Use
- Writing unit and integration tests in Go.
- Defining assertions (comparisons, errors, type assertions).
- Mocking internal services or database interfaces for unit tests.
- Implementing setup/teardown lifecycles for test environments using suites.
- Organizing and separating unit and integration test executions.

## Step-by-Step Instructions

1. **Verify Assertion Style**: Read `assets/assert-example.go.template` and apply `require` for immediate failures (e.g., error checks) and `assert` for subsequent value assertions.
2. **Handle Environment Isolation**: Use `t.Setenv(...)` instead of `os.Setenv` to ensure environment variables are localized and cleaned up automatically.
3. **Configure Build Tags and Makefile**: Read `references/integration-testing.md` to configure Go build tags for integration-only tests, and set up distinct test-unit and test-integration targets in the Makefile.

## Common Mistakes
- **Mutating Global Process Environment (TD-006)**: Modifying env vars with `os.Setenv` inside tests, which pollutes the global test runner state and breaks parallel execution. Use `t.Setenv` instead.
- **Polluting Unit Tests with Integration Logic**: Letting integration tests run as part of standard `go test ./...` without build tags, which breaks local developer pipelines if the external DB is off (BE-005).
- **Relying on Cached Test Results**: Forgetting `-count=1`, which allows Go to report success from cached results even if the database configuration has drifted or failed since the last run.
- **Failing to use `.AssertExpectations(t)` on mock**: Forgetting to verify that the mock methods were actually called. Always call this at the end of the test.
- **Using Assert instead of Require for err checks**: Using `assert.NoError` before checking fields of a returned struct pointer. If the pointer is nil, subsequent calls will panic. Always use `require.NoError` first.
