# Implementation Plan: OpenBench Initialization

## 1. Context Acquisition
- **Search & Map**: We have a newly created repository for `OpenBench`. The requirements dictate setting up a basic Go application structure tailored for RESTful API services. Additionally, it requires a centralized task runner (`Makefile`) to manage daily development workflows.
- **Dependency Audit**: The primary upstream dependencies will be:
    - Go runtime standard libraries.
    - `github.com/gofiber/fiber/v3` for the web framework.
    - `go.uber.org/zap` for structured logging.
    - `github.com/air-verse/air` for development-time hot reloading.

## 2. The Implementation Plan (Structure)

### A. Logical Requirements
- Establish a scoped Go module (`github.com/denden-dr/OpenBench`) for proper package and dependency resolution.
- Configure a task runner (`Makefile`) to standardize development commands (`build`, `run`, `fmt`, `tidy`).
- Configure a web server that listens on a specified port.
- Introduce a structured logger system using `zap` to capture application events.
- Implement a custom Fiber middleware utilizing `zap` to automatically log incoming HTTP access logs.
- Provide a dedicated `/health` endpoint to reflect the server's operational state.
- Embed an `.air.toml` configuration to monitor file changes and trigger automatic rebuilds for faster development iteration.

### B. Structural Strategy
- **File System Impact**:
    - `go.mod` / `go.sum`: Managed through Go module initialization.
    - `Makefile`: Script file containing all required shell execution phases.
    - `cmd/api/main.go`: Entry point for the compiled HTTP application.
    - `internal/handlers/health.go`: Component isolating the logic for the health check.
    - `internal/middleware/logger.go`: Custom middleware component generating HTTP access logs using `zap`.
    - `pkg/logger/logger.go`: Wrapper component encapsulating logging library logic and configuration.
    - `.air.toml`: Configuration file dictating the hot-reloading behavior rules (directories to ignore, binary paths, etc.).
- **Module Architecture**: 
    - The core application boots by initializing its local dependencies like the `logger` component first to ensure full observability. 
    - It then constructs a new Fiber application instance, injects the `zap` middleware, and maps its router paths. 
    - Finally, the server runs and starts listening.
- **Interface specs**: 
    - `health handler`: Accepts an HTTP context and immediately returns an HTTP `200 OK` structure to declare readiness.
    - `Makefile` specs:
        - `make run`: Starts the `air` hot reload system.
        - `make build`: Performs a standard `go build` outputting to a predefined `bin/` or `tmp/` target.
        - `make fmt`: Executes `go fmt ./...` to enforce standard Go formatting.
        - `make tidy`: Executes `go mod tidy` to clean and download all `go.mod` dependencies.

### C. Step-by-Step Logic
1. **Validation & Initialization**: 
   - Initialize the `go.mod` with the scoped package name `github.com/denden-dr/OpenBench`.
   - Install required dependencies.
2. **Configure Task Runner (Makefile)**: 
   - Generate the Makefile and populate it with `build`, `run`, `fmt`, and `tidy` targets.
   - Ensure the default state provides helper text or defaults to `run`.
3. **Setup Logger Utility**: 
   - Define a logger object mapping to `zap` standard capabilities initialized within `pkg/logger`.
   - Scaffold a Fiber middleware under `internal/middleware` that extracts request lifecycle data (status, time, path) and passes it to `zap`.
4. **Build Health Handler**: 
   - Define a handler method for the Fiber v3 router, returning a `{ status: "ok" }` or similar payload.
5. **Wire the Application Context**: 
   - Formulate `cmd/api/main.go` to construct the Fiber object, register the custom logging middleware, and attach the health handler at `/health`.
6. **Persist the Application Loop**: 
   - Initiate the Fiber `Listen` phase binding to a specified network port securely.
7. **Configure Hot Reloading**: 
   - Initialize `.air.toml` specifying the build command for `cmd/api/main.go` and executing it.

### D. Best Practice & Quality Guardrails
- **Task Runner Centralization**: Centralizing operations logically within the Makefile ensures consistent environments amongst any future contributors.
- **Error Handling**: Monitor the Fiber app bindings. If network assignment fails, emit a fatal signal to the logger and terminate gracefully.
- **Security**: Adopt the framework's baseline safe HTTP operational modes.
- **Observability**: Avoid `fmt.Print` or default `log` packages. All application-level messaging should leverage the customized `zap` configuration.

## 3. Verification Plan

### Automated Tests
- Validate basic Go compiler guarantees by running `make build`.
- Validate code linting guarantees by checking output from `make fmt`.

### Manual Verification
1. Run `make tidy` to confirm dependencies are securely tracked.
2. Run `make run` to invoke hot reloading via Air. Ensure no runtime warnings.
3. Access the web endpoint locally using `curl http://localhost:<port>/health` ensuring a valid response code.
4. Perform an ad-hoc change in the `health.go` string payload and ensure the terminal immediately acknowledges re-compilation and the new message propagates on the subsequent request.
