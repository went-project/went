package project

import "fmt"

func ReadmeTemplate(projectName string) string {
	return fmt.Sprintf(`# %[1]s

A REST API project scaffolded with [went](https://github.com/ekinatam/went).

> **For AI Agents:** See [CLAUDE.md](./CLAUDE.md) for agent guidelines, which points to [AGENTS.md](./AGENTS.md) for detailed conventions.

- [went](https://github.com/ekinatam/went) CLI (optional, for code generation)

## Getting Started

%[2]s%[2]s%[2]sbash
# 1. Copy environment file and fill in your values
cp .env.example .env

# 2. Install dependencies
go mod tidy

# 3. Run migrations
went migrate

# 4. Start the server
went run
%[2]s%[2]s%[2]s

The API will be available at %[2]shttp://localhost:8080%[2]s.
Swagger UI is served at %[2]shttp://localhost:8080/swagger/index.html%[2]s.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| %[2]sPORT%[2]s | %[2]s8080%[2]s | HTTP server port |
| %[2]sDB_CONNECTION%[2]s | %[2]ssqlite%[2]s | Database driver (%[2]ssqlite%[2]s / %[2]smysql%[2]s / %[2]spostgres%[2]s) |
| %[2]sDB_HOST%[2]s | %[2]s127.0.0.1%[2]s | Database host |
| %[2]sDB_PORT%[2]s | %[2]s3306%[2]s | Database port |
| %[2]sDB_NAME%[2]s | %[2]sdatabase%[2]s | Database name |
| %[2]sDB_USER%[2]s | — | Database user |
| %[2]sDB_PASSWORD%[2]s | — | Database password |
| %[2]sDB_STORAGE%[2]s | %[2]s./database.sqlite%[2]s | SQLite file path |
| %[2]sJWT_SECRET%[2]s | — | Secret key for JWT tokens |

If no env file is present, default DB settings fall back to SQLite.

## Code Generation

Use the %[2]swent%[2]s CLI to scaffold new resources:

%[2]s%[2]s%[2]sbash
went gen:model <Name>              # Skeleton related set
went gen:model <Name> -m           # Compatibility flag (skeleton related set already includes migration)
went gen:model <Name> -a           # Full related set
went gen:controller <Name>         # Skeleton related set
went gen:controller <Name> -a      # Full related set
went gen:router <Name>             # Skeleton related set
went gen:router <Name> -a          # Full related set
went gen:resource <Name>           # Skeleton related set
went gen:resource <Name> -a        # Full related set
went gen:migration <Name>          # Skeleton related set
went gen:migration <Name> -a       # Full related set
%[2]s%[2]s%[2]s

Default mode generates skeleton files so you can fill your domain details manually.
Use %[2]s--all%[2]s to generate fully populated templates.

Commands other than %[2]screate%[2]s require %[2]swentconfig.json%[2]s in the current directory.
If the file is missing, the CLI stops with an error.

Generated router files follow %[2]sroutes/<name>_router.go%[2]s naming in lowercase.

## Migrations

%[2]s%[2]s%[2]sbash
went migrate                       # Run all pending migrations
went migrate:rollback              # Roll back the last migration
went migrate:rollback --step 3     # Roll back the last 3 migrations
went migrate:fresh                 # Drop everything and re-run all migrations
%[2]s%[2]s%[2]s

## Project Structure

%[2]s%[2]s%[2]s
%[1]s/
├── main.go
├── go.mod
├── .env.example
├── database/
│   ├── migrations/
│   ├── models/
│   └── seeders/
├── http/
│   ├── controllers/
│   ├── middlewares/
│   ├── requests/
│   └── resources/
├── internal/
│   ├── config/
│   ├── responses/
│   └── providers/
└── routes/
%[2]s%[2]s%[2]s

## API Documentation

After running %[2]sswag init%[2]s and starting the server, open:

%[2]s%[2]s%[2]s
http://localhost:8080/swagger/index.html
%[2]s%[2]s%[2]s

## Resource Query Layer

Generated resources under %[2]shttp/resources/%[2]s include:

- %[2]sXResource%[2]s: single-item API payload
- %[2]sXCollection%[2]s: Laravel-style paginated payload with %[2]sdata%[2]s and %[2]smeta%[2]s
- %[2]sXQuery%[2]s: DB read operations (%[2]sPaginate%[2]s, %[2]sFind%[2]s) that controllers delegate to

Shared pagination metadata lives in %[2]shttp/resources/pagination.go%[2]s.

## Request DTO Layer

Generated request DTOs under %[2]shttp/requests/%[2]s include:

- %[2]sXPayload%[2]s: create payload
- %[2]sXUpdatePayload%[2]s: partial update payload

Controllers import these request types instead of defining payload structs inline.

## Global Error Responses

Generated helper file %[2]sinternal/responses/error_response.go%[2]s provides shared API error primitives:

- %[2]sErrorBody%[2]s: standard error payload schema
- %[2]sJSONError(...)%[2]s: common JSON error response helper
- %[2]sJSONErrorWithDetails(...)%[2]s: JSON error response with optional details

Controllers are generated to use these helpers and Swagger failure responses reference %[2]sresponses.ErrorBody%[2]s.
`, projectName, "`")
}
