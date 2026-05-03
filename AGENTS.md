# AGENTS.md

## Purpose

This file documents how coding agents should work in this repository.

> See [CLAUDE.md](./CLAUDE.md) for a quick introduction to these guidelines.

## Repository Rules

- Keep changes minimal and focused to the requested feature.
- Preserve existing command naming style (`gen:*`).
- Prefer templates under `internal/templates` for scaffold content.
- Prefer handlers under `internal/handlers` for generation logic.
- Do not hardcode absolute paths.
- Keep generated file names lowercase where current conventions use lowercase.
- Each `gen:*` command lives in its own file under `cmd/` (`gen_router.go`, `gen_model.go`, `gen_controller.go`, `gen_migration.go`).

## CLI Conventions

- Main binary name from `Makefile`: `went`
- Root commands:
  - `create [name]`
  - `version`
  - `gen:router [name]`
  - `gen:model [name]`
  - `gen:controller [name]`
  - `gen:migration [name]`
  - `gen:resource [name]`
  - `migrate`
  - `migrate:rollback` (`--step/-s N`, default 1)
  - `migrate:fresh`
- Each command is registered in `cmd/root.go` via `Root.AddCommand(...)`.
- `cmd/gen.go` is kept as an empty package stub; all command vars live in dedicated files.
- `create` disindaki tum komutlar, calisma dizininde `wentconfig.json` oldugunu dogrulamalidir; dosya yoksa komut hata ile durmalidir.

## Generation Flow Expectations

- `create [name]` should scaffold a runnable project layout.
- `create` runs `go mod init`, `go mod tidy`, and `swag init` automatically. Do NOT add `go get -u` — it breaks on Go < 1.25.
- All sub-commands inside `CreateProject` must use the `runCommand` helper so errors include stdout/stderr output.
- New generation capabilities should be exposed both as:
  - template functions in `internal/templates`
  - dedicated handlers in `internal/handlers`
  - user-facing Cobra commands in `cmd` (one file per command)
- Migration generation should produce both `.up.sql` and `.down.sql` files.
- `migrate`, `migrate:rollback`, and `migrate:fresh` live in `cmd/migrate.go`, `cmd/migrate_rollback.go`, `cmd/migrate_fresh.go` respectively. Logic lives in `internal/handlers/Migrations.go`.
- The migration tracking table is named `wentmigrations`. Do NOT rename it to `schema_migrations` or anything else.
- `migrate:rollback --step N` rolls back the last N applied migrations in reverse order.
- `migrate:fresh` runs all down files (ignoring errors), drops `wentmigrations`, then re-runs all up files.
- Every `gen:*` command supports `-a/--all` to generate all related templates in full form.
- Without `--all`, `gen:*` commands generate all related templates in skeleton form and leave implementation details to the user.
- `gen:model` keeps `-m/--migration` for compatibility; in skeleton related generation mode migration creation is already included.
- `gen:resource` lives in `cmd/gen_resource.go`; handler lives in `internal/handlers/CreateResource.go`.
- `CreateResource` must ensure `http/resources/pagination.go` exists, then generate `<name>_resource.go`.
- `CreateRouter` generated output file must be `routes/<lowercase_name>_router.go`.
- `create` must generate `internal/responses/error_response.go` via `CreateGlobalErrorResponse(...)`.

## Template Conventions

### Controller Template (`internal/templates/http/controller_template.go`)
- Uses shared `ParseID` from `http/controllers/helpers.go` for route ID parsing.
- Keeps `findByID` helper method for model lookup reuse.
- Does NOT generate local `XPayload` or `XResponse` structs.
- Uses request DTOs from `http/requests` (`requests.XPayload`, `requests.XUpdatePayload`).
- Uses response DTOs from `http/resources` (`resources.XResource`, `resources.XCollection`).
- All CRUD methods include full Swagger annotations: `@Summary`, `@Description`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`.
- `@Router` paths use lowercase resource name (e.g. `User` → `/user`).
- `Update` uses `map[string]interface{}` + `DB.Model(&item).Updates(updates)` — do NOT hardcode field assignments.
- `GetAll` and `GetByID` should delegate read queries to `http/resources` query builders (`resources.NewXQuery(...)`) instead of duplicating query logic in controller methods.
- Error responses should use `internal/responses` helpers (`responses.JSONError`, `responses.JSONErrorWithDetails`) instead of `http.Error`.
- `@Failure` annotations should reference `responses.ErrorBody` for standard error payload docs.

### Request Template (`internal/templates/http/request_template.go`)
- Must generate `XPayload` and `XUpdatePayload` in `http/requests`.
- Must provide both full and skeleton template variants to support generation modes.

### Controller Helper Template (`internal/templates/http/controller_helper_template.go`)
- Must generate shared `ParseID(...)` in `http/controllers/helpers.go`.
- `CreateControllerHelper(...)` should ensure this file exists for generated controllers.

### Resource Template (`internal/templates/http/resource_template.go`)
- Must generate `XResource`, `XCollection`, and `XQuery` types.
- `XQuery` owns DB read operations (`Paginate`, `Find`) and returns transformed resource payloads.
- `Paginate` should build Laravel-style envelope: `{ data: [...], meta: {...} }` using shared pagination metadata.

### Pagination Template (`internal/templates/http/pagination_template.go`)
- Provides shared `PaginationMeta` and `BuildMeta(...)` used by resource collections.
- Generated as `http/resources/pagination.go` once per project.

### Response Template (`internal/templates/responses/error_template.go`)
- Generates `internal/responses/error_response.go` with shared API error response primitives.
- Must expose `ErrorBody`, `JSONError(...)`, and `JSONErrorWithDetails(...)`.

### Config Template (`internal/templates/config/config_template.go`)
- Reads both `DB_NAME` and `DB_STORAGE` env vars.
- If no env file is found, defaults should still lead to SQLite usage.
- For SQLite dialect: if `DB_NAME` is the default (`database`), it falls back to `DB_STORAGE`. This prevents the generated `database/` folder from being opened as a file.
- Default SQLite path: `./database.sqlite`.

### Database Provider Template (`internal/templates/database/database_provider_template.go`)
- Panic message includes the underlying error: `fmt.Sprintf("Failed to connect database: %v", err)`.

### Main Template (`internal/templates/app/main_template.go`)
- Imports `_ "<appname>/docs"` for Swagger docs side-effect registration.
- All import strings must be properly terminated (no missing closing quotes).

## Editing Guidance

- Update docs when commands or scaffold outputs change.
- Validate with diagnostics after edits.
- Prefer deterministic template output for reproducible scaffolding.
- After editing any template, verify the CLI still builds with `go build ./...`.
