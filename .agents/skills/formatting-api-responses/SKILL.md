---
name: formatting-api-responses
description: Use when building or modifying HTTP handlers in the Go backend to ensure standardized responses and validation. Do not use for non-Go backend frameworks or frontend API client formatting.
version: 1.0.0
---

# Formatting API Responses & Validation

## Overview
Every backend HTTP handler must use the standardized `response` package to format both success and error responses. Requests should be validated using structural tags and the `validator` helper package. Handlers must be implemented as methods on a struct rather than closure-returning functions.

## When to Use
- Writing or refactoring Go HTTP handlers in Fiber.
- Creating or editing request and response DTO schemas.
- Adding input validation rules to endpoint parameters.

## Step-by-Step Instructions

1. **Organize DTOs**: Create DTO schemas inside a local `dto.go` file (e.g. `internal/auth/dto.go`) to prevent import cycle locks. Read `assets/dto.go.template` for structure and validation tags.
2. **Implement Handler Structure**: Define handlers as methods of a Handler struct. Read `assets/handler.go.template` to implement the struct, constructor, body parser, and validator execution.
3. **Register routes**: Read `references/routing-example.md` to map struct methods directly to Fiber's router routes.
4. **Standardize responses**: Use the standard `response` package helper methods (`response.JSON` and `response.Error`) to ensure JSON responses are properly enveloped.

## Common Mistakes to Avoid
- **Closure Handlers**: Returning `fiber.Handler` from a custom function. Use a `Handler` struct constructor and register methods instead.
- **Direct fiber.Map / Raw JSON Errors**: Bypassing the `response` package. Always use the helper methods `response.JSON` or `response.Error` to guarantee format consistency.
- **Global DTOs**: Placing endpoint-specific request/response DTOs in a global space, risking import cycle locks. Keep them locally scoped in the feature package's `dto.go`.
