---
name: initializing-go-fiber-api
description: Use when initializing a Go web application using the Fiber framework, configuring security middleware (CORS, Recover, Helmet), and implementing graceful shutdown. Do not use for Gin, Chi, Echo, or other non-Fiber web frameworks.
version: 1.1.0
---

# Initializing Go Fiber API

## Overview
Go Fiber is an Express-inspired web framework for Go. When setting up a Fiber API, security middleware, fail-fast server startup, and graceful shutdown must be implemented from the start.

## When to Use
- Initializing a new Go web server using Fiber.
- Hardening an existing Fiber application with security headers, recovery logic, or CORS.
- Adding graceful shutdown handling to a Go service.

## Step-by-Step Instructions

1. **Initialize Fiber App**: Read `assets/server.go.template` to understand the standard Fiber initialization flow, establishing appropriate read/write timeouts to protect against slowloris attacks.
2. **Register Security Middlewares**: Configure `recover.New()` first to capture unhandled panics, and configure `cors.New(...)` with strict, explicit domain origins. When using cookie-based authentication, ensure CORS config sets `AllowCredentials: true` and lists explicit origins, as browsers reject cookie transfer on wildcard `*` origins.
3. **Implement Fail-Fast Port Binding**: Instantiate a buffered error channel of size 1 (e.g. `make(chan error, 1)`), run `app.Listen` inside a goroutine, and write any returned listener errors to the error channel.
4. **Implement Graceful Shutdown**: Establish an OS signal interrupt channel listening for `SIGINT` and `SIGTERM`. Block using a `select` statement that monitors both the startup error channel and the OS signal channel. If an OS signal is caught, trigger `app.Shutdown()` to allow active requests to finish cleanly.

## Common Mistakes
- **Hanging Startup**: Running `app.Listen` inside a goroutine without returning its errors, which leaves the main process alive but unable to serve traffic when port bind fails (BE-001).
- **No Recover Middleware**: Not using the `recover` middleware. Any unhandled panic in a route handler will crash the entire server.
- **Permissive CORS**: Using `AllowOrigins: "*"` in production. Be restrictive and define exact domain origins.
- **No Read/Write Timeouts**: Leaving timeouts undefined in `fiber.Config`. This leaves the server vulnerable to slowloris attacks.
- **Abrupt Termination**: Killing the process with `os.Exit` or allowing the main function to end immediately without calling `app.Shutdown()`, which terminates in-flight client requests mid-transaction.
