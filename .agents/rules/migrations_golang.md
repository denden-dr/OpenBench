# Golang Migrate Rules

When creating or managing database migrations, you must STRICTLY adhere to the following rules:

## 1. Migration Generation
- Migrations MUST NOT be created manually.
- You MUST use the `migrate` CLI tool (golang-migrate) to generate migration files.
- Prefer using the `make migrate-create` command if available in the project's `Makefile`.

## 2. File Naming Convention
- Migrations must use the sequential numbering format (e.g., `000001_name.up.sql`).
- Every migration must have both an `.up.sql` (for applying changes) and a `.down.sql` (for rolling back changes).

## 3. Makefile Integration
- The `Makefile` must contain a `migrate-create` target that automates the generation of these files.
- Example pattern for `Makefile`:
  ```makefile
  migrate-create:
      @read -p "Enter migration name: " name; \
      migrate create -ext sql -dir migrations -seq $$name
  ```

## Why this is required
This ensures consistency in migration versions, prevents accidental manual errors, and provides a standard way to rollback changes in any environment.
