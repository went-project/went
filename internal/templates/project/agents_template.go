package project

import "fmt"

func AgentsTemplate(projectName string) string {
	return fmt.Sprintf(`# AGENTS.md

## Purpose

This file documents how coding agents should work in the **%[1]s** project.

> For a quick introduction, see [CLAUDE.md](./CLAUDE.md).
├── main.go                        # Entry point, wires router and starts server
├── go.mod
├── .env.example                   # Environment variable reference
├── .gitignore
├── database/
│   ├── migrations/                # SQL migration files (*.up.sql / *.down.sql)
│   ├── models/                    # GORM model structs
│   └── seeders/                   # Database seeders
├── http/
│   ├── controllers/               # Request handlers (one file per resource)
│   ├── middlewares/               # auth.go, cors.go
│   ├── requests/                  # Request validation structs
│   └── resources/                 # Response transformation structs
├── internal/
│   ├── config/config.go           # Env-based config loader
│   ├── responses/error_response.go # Shared API error response helpers
│   └── providers/database_provider.go  # GORM DB initializer
└── routes/
    ├── main_router.go             # Gin engine setup + middleware mounting
    └── *_router.go                # Per-resource route groups
%[2]s%[2]s%[2]s

## Conventions

- **Models** live in %[2]sdatabase/models/%[2]s. Embed %[2]sgorm.Model%[2]s for standard fields.
- **Controllers** live in %[2]shttp/controllers/%[2]s. One file per resource, named %[2]s<resource>_controller.go%[2]s.
- Shared controller helpers live in %[2]shttp/controllers/helpers.go%[2]s (e.g. %[2]sParseID%[2]s).
- **Requests** live in %[2]shttp/requests/%[2]s. One file per resource, named %[2]s<resource>_request.go%[2]s.
- **Resources** live in %[2]shttp/resources/%[2]s. One file per resource, named %[2]s<resource>_resource.go%[2]s.
- **Routers** live in %[2]sroutes/%[2]s. One file per resource, named %[2]s<resource>_router.go%[2]s.
- **Migrations** live in %[2]sdatabase/migrations/%[2]s as numbered pairs: %[2]sNNNNNN_<desc>.up.sql%[2]s / %[2]sNNNNNN_<desc>.down.sql%[2]s.
- Keep file names lowercase.
- Commands other than %[2]screate%[2]s must be run in a directory that contains %[2]swentconfig.json%[2]s.
- Do not declare payload/response DTOs in controllers; use %[2]shttp/requests%[2]s and %[2]shttp/resources%[2]s types.
- In default generation mode, fill skeleton implementations manually.

## Generation Modes

- Default %[2]sgen:* <Name>%[2]s commands generate related files in skeleton form.
- %[2]sgen:* <Name> --all%[2]s generates related files in full form.
- %[2]sgen:model -m%[2]s is preserved for compatibility; skeleton related generation already includes migration files.

## Environment Variables

See %[2]s.env.example%[2]s for all required variables. Copy it to %[2]s.env%[2]s before running.
If no env file is found, configuration defaults still fall back to SQLite.

## Key Commands

| Command | Description |
|---------|-------------|
| %[2]sgo run main.go%[2]s | Start the development server |
| %[2]swent run%[2]s | Start the project with hot reload |
| %[2]sswag init%[2]s | Regenerate Swagger docs |
| %[2]swent gen:model <Name>%[2]s | Scaffold a new model |
| %[2]swent gen:resource <Name>%[2]s | Scaffold a new resource |
| %[2]swent gen:controller <Name>%[2]s | Scaffold a new controller |
| %[2]swent gen:router <Name>%[2]s | Scaffold a new router |
| %[2]swent gen:migration <Name>%[2]s | Scaffold a new migration |
| %[2]swent migrate%[2]s | Run pending migrations |
| %[2]swent migrate:rollback%[2]s | Roll back the last migration |
| %[2]swent migrate:fresh%[2]s | Drop all tables and re-run migrations |

## Editing Guidance

- Register new routes in the corresponding %[2]sroutes/*_router.go%[2]s file and mount the group in %[2]sroutes/main_router.go%[2]s.
- After adding or changing model structs, create a new migration pair and run %[2]swent migrate%[2]s.
- Keep read/query logic in %[2]shttp/resources/*_resource.go%[2]s via %[2]sXQuery%[2]s methods; keep controllers thin.
- Use %[2]sPaginate(page, perPage)%[2]s in resources for Laravel-style %[2]s{data, meta}%[2]s responses.
- Use %[2]sinternal/responses%[2]s helpers in controllers for error output (%[2]sJSONError%[2]s / %[2]sJSONErrorWithDetails%[2]s).
- Keep Swagger failure schemas aligned with %[2]sresponses.ErrorBody%[2]s.
- After changing any controller annotations, run %[2]sswag init%[2]s to update docs.
- Do not modify generated migration files that have already been applied.
`, projectName, "`")
}
